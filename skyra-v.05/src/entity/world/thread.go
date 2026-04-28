// Thread tracks an ongoing conversation in a system world. Threads are
// physics — they live on the world, not on the being.
package world

import (
	"crypto/rand"
	"fmt"
)

type Thread struct {
	id     string
	About  string
	Active bool
}

func NewThreadID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
