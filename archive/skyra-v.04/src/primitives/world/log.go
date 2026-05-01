package world

import (
	"fmt"
	"os"
	"path/filepath"

	"skyra-v04/src/primitives/entity"
	"skyra-v04/src/primitives/thread"
)

const logDir = "logs"

func logPresent(beingID, present string) {
	dir := filepath.Join(logDir, beingID)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return
	}
	_ = os.WriteFile(filepath.Join(dir, "present.txt"), []byte(present), 0644)
}

func logThreadState(beingID, stage string, t thread.Thread) {
	dir := filepath.Join(logDir, beingID)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return
	}
	f, err := os.OpenFile(filepath.Join(dir, "thread-state.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	fmt.Fprintf(f, "[%s] thread=%s\n", stage, t.ID())
	for key, ex := range t.ExchangeMap() {
		fmt.Fprintf(f, "  E{%s, %s}: parent=%s active=%v entries=%d\n", key.A, key.B, ex.Parent, ex.Active, len(ex.Relations))
	}
	fmt.Fprintln(f, "")
}

func logDrop(beingID, line, reason, present string) {
	dir := filepath.Join(logDir, beingID)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return
	}
	f, err := os.OpenFile(filepath.Join(dir, "dropped.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		fmt.Fprintf(f, "REASON: %s\nLINE:   %s\n\n", reason, line)
		f.Close()
	}
	body := fmt.Sprintf("=== last error ===\nreason: %s\nline:   %s\n\n=== present at time of error ===\n%s", reason, line, present)
	_ = os.WriteFile(filepath.Join(dir, "last-logged-error.txt"), []byte(body), 0644)
}

func traceRelation(stage string, r entity.Relation) {
	fmt.Printf("trace: %s origin=%s target=%s thread=%s impulse=%q\n", stage, r.Origin, r.ID, r.ThreadID, r.Impulse)
}
