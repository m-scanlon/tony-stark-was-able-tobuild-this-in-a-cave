package kernel

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

var ErrValkeyNil = errors.New("valkey: nil response")

// ValkeyError is a RESP protocol error returned by the server.
type ValkeyError struct {
	Message string
}

func (e *ValkeyError) Error() string {
	return "valkey: " + e.Message
}

// ValkeyClient is a tiny zero-dependency RESP client.
// It opens one connection per command for simplicity and isolation.
type ValkeyClient struct {
	addr        string
	username    string
	password    string
	dialTimeout time.Duration
	rwTimeout   time.Duration
}

func NewValkeyClient(addr, username, password string) *ValkeyClient {
	return &ValkeyClient{
		addr:        strings.TrimSpace(addr),
		username:    strings.TrimSpace(username),
		password:    password,
		dialTimeout: 5 * time.Second,
		rwTimeout:   5 * time.Second,
	}
}

func (c *ValkeyClient) Ping(ctx context.Context) error {
	resp, err := c.exec(ctx, []string{"PING"})
	if err != nil {
		return err
	}
	pong, ok := resp.(string)
	if !ok || strings.ToUpper(pong) != "PONG" {
		return fmt.Errorf("unexpected PING response: %v", resp)
	}
	return nil
}

func (c *ValkeyClient) Get(ctx context.Context, key string) (string, error) {
	resp, err := c.exec(ctx, []string{"GET", key})
	if err != nil {
		return "", err
	}
	value, ok := resp.(string)
	if !ok {
		return "", fmt.Errorf("GET %q: unexpected response type %T", key, resp)
	}
	return value, nil
}

func (c *ValkeyClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	resp, err := c.exec(ctx, []string{"HGETALL", key})
	if err != nil {
		return nil, err
	}

	items, ok := resp.([]any)
	if !ok {
		return nil, fmt.Errorf("HGETALL %q: unexpected response type %T", key, resp)
	}
	if len(items)%2 != 0 {
		return nil, fmt.Errorf("HGETALL %q: invalid field/value pair count %d", key, len(items))
	}

	fields := make(map[string]string, len(items)/2)
	for i := 0; i < len(items); i += 2 {
		k, ok := items[i].(string)
		if !ok {
			return nil, fmt.Errorf("HGETALL %q: field type %T is not string", key, items[i])
		}
		v, ok := items[i+1].(string)
		if !ok {
			return nil, fmt.Errorf("HGETALL %q: value type %T is not string", key, items[i+1])
		}
		fields[k] = v
	}
	return fields, nil
}

func (c *ValkeyClient) Exists(ctx context.Context, key string) (bool, error) {
	resp, err := c.exec(ctx, []string{"EXISTS", key})
	if err != nil {
		return false, err
	}
	count, ok := resp.(int64)
	if !ok {
		return false, fmt.Errorf("EXISTS %q: unexpected response type %T", key, resp)
	}
	return count > 0, nil
}

func (c *ValkeyClient) Set(ctx context.Context, key, value string) error {
	resp, err := c.exec(ctx, []string{"SET", key, value})
	if err != nil {
		return err
	}
	ok, okType := resp.(string)
	if !okType || strings.ToUpper(ok) != "OK" {
		return fmt.Errorf("SET %q: unexpected response %v", key, resp)
	}
	return nil
}

func (c *ValkeyClient) exec(ctx context.Context, args []string) (any, error) {
	if c == nil || c.addr == "" {
		return nil, fmt.Errorf("valkey client is not configured")
	}

	dialer := &net.Dialer{Timeout: c.dialTimeout}
	conn, err := dialer.DialContext(ctx, "tcp", c.addr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if deadline, ok := ctx.Deadline(); ok {
		if err := conn.SetDeadline(deadline); err != nil {
			return nil, err
		}
	} else {
		if err := conn.SetDeadline(time.Now().Add(c.rwTimeout)); err != nil {
			return nil, err
		}
	}

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	if c.password != "" {
		authArgs := []string{"AUTH", c.password}
		if c.username != "" {
			authArgs = []string{"AUTH", c.username, c.password}
		}
		if err := writeCommand(writer, authArgs...); err != nil {
			return nil, err
		}
		if err := writer.Flush(); err != nil {
			return nil, err
		}
		if _, err := readRESP(reader); err != nil {
			return nil, err
		}
	}

	if err := writeCommand(writer, args...); err != nil {
		return nil, err
	}
	if err := writer.Flush(); err != nil {
		return nil, err
	}

	return readRESP(reader)
}

func writeCommand(w io.Writer, args ...string) error {
	if _, err := fmt.Fprintf(w, "*%d\r\n", len(args)); err != nil {
		return err
	}
	for _, arg := range args {
		if _, err := fmt.Fprintf(w, "$%d\r\n%s\r\n", len(arg), arg); err != nil {
			return err
		}
	}
	return nil
}

func readRESP(r *bufio.Reader) (any, error) {
	prefix, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	switch prefix {
	case '+':
		line, err := readLine(r)
		if err != nil {
			return nil, err
		}
		return line, nil
	case '-':
		line, err := readLine(r)
		if err != nil {
			return nil, err
		}
		return nil, &ValkeyError{Message: line}
	case ':':
		line, err := readLine(r)
		if err != nil {
			return nil, err
		}
		value, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			return nil, err
		}
		return value, nil
	case '$':
		line, err := readLine(r)
		if err != nil {
			return nil, err
		}
		size, err := strconv.Atoi(line)
		if err != nil {
			return nil, err
		}
		if size == -1 {
			return nil, ErrValkeyNil
		}
		buf := make([]byte, size+2)
		if _, err := io.ReadFull(r, buf); err != nil {
			return nil, err
		}
		return string(buf[:size]), nil
	case '*':
		line, err := readLine(r)
		if err != nil {
			return nil, err
		}
		size, err := strconv.Atoi(line)
		if err != nil {
			return nil, err
		}
		if size == -1 {
			return nil, ErrValkeyNil
		}
		items := make([]any, 0, size)
		for i := 0; i < size; i++ {
			item, err := readRESP(r)
			if err != nil {
				return nil, err
			}
			items = append(items, item)
		}
		return items, nil
	default:
		return nil, fmt.Errorf("valkey: unknown RESP prefix %q", prefix)
	}
}

func readLine(r *bufio.Reader) (string, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	line = strings.TrimSuffix(line, "\n")
	line = strings.TrimSuffix(line, "\r")
	return line, nil
}
