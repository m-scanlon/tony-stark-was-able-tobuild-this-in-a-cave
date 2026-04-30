package debug

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	file *os.File
	mu   sync.Mutex
)

func Init(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	name := time.Now().Format("2006-01-02_15-04-05") + ".log"
	f, err := os.Create(filepath.Join(dir, name))
	if err != nil {
		return err
	}
	file = f
	return nil
}

func Log(args ...any) {
	mu.Lock()
	defer mu.Unlock()
	if file == nil {
		return
	}
	fmt.Fprintln(file, args...)
}

func Close() {
	mu.Lock()
	defer mu.Unlock()
	if file != nil {
		file.Close()
	}
}
