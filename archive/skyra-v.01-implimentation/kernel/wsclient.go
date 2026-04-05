package kernel

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"
	"sync"
	"time"
)

type wsClient struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
	mu     sync.Mutex
}

func dialWebSocket(ctx context.Context, rawURL string) (*wsClient, error) {
	u, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return nil, fmt.Errorf("parse websocket url: %w", err)
	}
	if u.Scheme != "ws" && u.Scheme != "wss" {
		return nil, fmt.Errorf("unsupported websocket scheme %q", u.Scheme)
	}

	host := u.Host
	if !strings.Contains(host, ":") {
		if u.Scheme == "wss" {
			host += ":443"
		} else {
			host += ":80"
		}
	}

	var conn net.Conn
	dialer := &net.Dialer{Timeout: 10 * time.Second}
	switch u.Scheme {
	case "wss":
		conn, err = tls.DialWithDialer(dialer, "tcp", host, &tls.Config{
			ServerName: strings.Split(u.Host, ":")[0],
			MinVersion: tls.VersionTLS12,
		})
	default:
		conn, err = dialer.DialContext(ctx, "tcp", host)
	}
	if err != nil {
		return nil, err
	}

	keyRaw := make([]byte, 16)
	if _, err := rand.Read(keyRaw); err != nil {
		conn.Close()
		return nil, err
	}
	key := base64.StdEncoding.EncodeToString(keyRaw)

	path := u.RequestURI()
	if path == "" {
		path = "/"
	}

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	headers := []string{
		fmt.Sprintf("GET %s HTTP/1.1", path),
		fmt.Sprintf("Host: %s", u.Host),
		"Upgrade: websocket",
		"Connection: Upgrade",
		fmt.Sprintf("Sec-WebSocket-Key: %s", key),
		"Sec-WebSocket-Version: 13",
		"",
		"",
	}
	if _, err := writer.WriteString(strings.Join(headers, "\r\n")); err != nil {
		conn.Close()
		return nil, err
	}
	if err := writer.Flush(); err != nil {
		conn.Close()
		return nil, err
	}

	tp := textproto.NewReader(reader)
	statusLine, err := tp.ReadLine()
	if err != nil {
		conn.Close()
		return nil, err
	}
	if !strings.Contains(statusLine, "101") {
		conn.Close()
		return nil, fmt.Errorf("websocket handshake failed: %s", statusLine)
	}
	mimeHeader, err := tp.ReadMIMEHeader()
	if err != nil {
		conn.Close()
		return nil, err
	}
	accept := mimeHeader.Get("Sec-WebSocket-Accept")
	expected := computeWebSocketAccept(key)
	if accept != expected {
		conn.Close()
		return nil, errors.New("websocket handshake accept mismatch")
	}

	return &wsClient{
		conn:   conn,
		reader: reader,
		writer: writer,
	}, nil
}

func computeWebSocketAccept(key string) string {
	sum := sha1.Sum([]byte(key + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
	return base64.StdEncoding.EncodeToString(sum[:])
}

func (c *wsClient) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func (c *wsClient) WriteText(payload string) error {
	return c.writeFrame(0x1, []byte(payload))
}

func (c *wsClient) WritePong(payload []byte) error {
	return c.writeFrame(0xA, payload)
}

func (c *wsClient) writeFrame(opcode byte, payload []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	first := byte(0x80 | opcode)
	if err := c.writer.WriteByte(first); err != nil {
		return err
	}

	maskKey := make([]byte, 4)
	if _, err := rand.Read(maskKey); err != nil {
		return err
	}

	payloadLen := len(payload)
	switch {
	case payloadLen < 126:
		if err := c.writer.WriteByte(byte(payloadLen) | 0x80); err != nil {
			return err
		}
	case payloadLen <= 65535:
		if err := c.writer.WriteByte(126 | 0x80); err != nil {
			return err
		}
		if err := binary.Write(c.writer, binary.BigEndian, uint16(payloadLen)); err != nil {
			return err
		}
	default:
		if err := c.writer.WriteByte(127 | 0x80); err != nil {
			return err
		}
		if err := binary.Write(c.writer, binary.BigEndian, uint64(payloadLen)); err != nil {
			return err
		}
	}

	if _, err := c.writer.Write(maskKey); err != nil {
		return err
	}

	masked := make([]byte, len(payload))
	for i := range payload {
		masked[i] = payload[i] ^ maskKey[i%4]
	}
	if _, err := c.writer.Write(masked); err != nil {
		return err
	}
	return c.writer.Flush()
}

func (c *wsClient) ReadText(ctx context.Context) (string, error) {
	for {
		if deadline, ok := ctx.Deadline(); ok {
			if err := c.conn.SetReadDeadline(deadline); err != nil {
				return "", err
			}
		}

		payload, opcode, err := c.readFrame()
		if err != nil {
			return "", err
		}

		switch opcode {
		case 0x1:
			return string(payload), nil
		case 0x8:
			return "", io.EOF
		case 0x9:
			if err := c.WritePong(payload); err != nil {
				return "", err
			}
		case 0xA:
		default:
		}
	}
}

func (c *wsClient) readFrame() ([]byte, byte, error) {
	first, err := c.reader.ReadByte()
	if err != nil {
		return nil, 0, err
	}
	second, err := c.reader.ReadByte()
	if err != nil {
		return nil, 0, err
	}

	opcode := first & 0x0F
	masked := second&0x80 != 0
	payloadLen := int(second & 0x7F)

	switch payloadLen {
	case 126:
		var n uint16
		if err := binary.Read(c.reader, binary.BigEndian, &n); err != nil {
			return nil, 0, err
		}
		payloadLen = int(n)
	case 127:
		var n uint64
		if err := binary.Read(c.reader, binary.BigEndian, &n); err != nil {
			return nil, 0, err
		}
		payloadLen = int(n)
	}

	var maskKey [4]byte
	if masked {
		if _, err := io.ReadFull(c.reader, maskKey[:]); err != nil {
			return nil, 0, err
		}
	}

	payload := make([]byte, payloadLen)
	if _, err := io.ReadFull(c.reader, payload); err != nil {
		return nil, 0, err
	}

	if masked {
		for i := range payload {
			payload[i] ^= maskKey[i%4]
		}
	}

	return payload, opcode, nil
}

func httpURLFromWS(raw string) (string, error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return "", err
	}
	switch u.Scheme {
	case "ws":
		u.Scheme = "http"
	case "wss":
		u.Scheme = "https"
	default:
		return "", fmt.Errorf("unsupported websocket scheme %q", u.Scheme)
	}
	u.Path = ""
	u.RawPath = ""
	u.RawQuery = ""
	u.Fragment = ""
	return strings.TrimRight(u.String(), "/"), nil
}

func checkHTTPReady(ctx context.Context, raw string) error {
	base, err := httpURLFromWS(raw)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, base+"/readyz", nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("gateway readyz returned %d", resp.StatusCode)
	}
	return nil
}
