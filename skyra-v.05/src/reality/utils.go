package reality

import "crypto/rand"

func NewID() string {
	b := make([]byte, 8)
	rand.Read(b)
	const hex = "0123456789abcdef"
	s := make([]byte, 16)
	for i, v := range b {
		s[i*2] = hex[v>>4]
		s[i*2+1] = hex[v&0x0f]
	}
	return string(s)
}
