package inference

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func CallClaude(system, present string) (string, error) {
	prompt := system + "\n\n" + present

	cmd := exec.Command("claude", "-p", "--output-format", "text")
	cmd.Stdin = strings.NewReader(prompt)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	done := make(chan error, 1)
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("claude start: %w", err)
	}
	go func() { done <- cmd.Wait() }()

	select {
	case err := <-done:
		if err != nil {
			errMsg := strings.TrimSpace(stderr.String())
			if errMsg == "" {
				errMsg = err.Error()
			}
			return "", fmt.Errorf("claude: %s", errMsg)
		}
		return strings.TrimSpace(stdout.String()), nil
	case <-time.After(120 * time.Second):
		cmd.Process.Kill()
		return "", fmt.Errorf("claude: timeout")
	}
}
