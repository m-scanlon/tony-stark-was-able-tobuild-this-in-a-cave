package medium

import (
	"fmt"
	"io"
	"os/exec"
	"strings"

	"skyra-v04/src/primitives/entity"
)

// execMedium returns a Medium bound to a specific binary path. Each call spawns
// the binary fresh, writes the incoming relation as a protocol line to its
// stdin, reads the response protocol from stdout, and returns it.
func execMedium(binPath string) Medium {
	return func(_ string, r entity.Relation) (string, error) {
		cmd := exec.Command(binPath)
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return "", fmt.Errorf("exec(%s): stdin pipe: %w", binPath, err)
		}
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return "", fmt.Errorf("exec(%s): stdout pipe: %w", binPath, err)
		}
		if err := cmd.Start(); err != nil {
			return "", fmt.Errorf("exec(%s): start: %w", binPath, err)
		}

		// Write the incoming relation as a protocol line to the binary's stdin.
		go func() {
			defer stdin.Close()
			fmt.Fprintf(stdin, "%s %s\n", r.ID, r.Impulse)
		}()

		out, err := io.ReadAll(stdout)
		if err != nil {
			return "", fmt.Errorf("exec(%s): read stdout: %w", binPath, err)
		}
		if err := cmd.Wait(); err != nil {
			return "", fmt.Errorf("exec(%s): wait: %w", binPath, err)
		}
		return strings.TrimSpace(string(out)), nil
	}
}
