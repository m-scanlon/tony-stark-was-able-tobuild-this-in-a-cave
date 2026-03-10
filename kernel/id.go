package kernel

import (
	"crypto/rand"
	"encoding/hex"
)

func newID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic("kernel: failed to generate ID: " + err.Error())
	}
	return hex.EncodeToString(b)
}
