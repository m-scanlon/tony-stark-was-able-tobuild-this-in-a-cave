package medium

import (
	"fmt"
	"os/exec"
	"strings"

	"skyra-v04/src/primitives/entity"
)

func init() {
	Register("shell", shell)
}

func shell(_ string, r entity.Relation) (string, error) {
	cmd := strings.TrimSpace(r.Impulse)
	if cmd == "" {
		return "", nil
	}
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	response := string(out)
	if err != nil {
		response = "error: " + err.Error() + "\n" + response
	}
	safe := strings.ReplaceAll(response, "|", "│")
	safe = strings.ReplaceAll(safe, "\n", " ↵ ")
	safe = strings.TrimSpace(safe)
	if safe == "" {
		safe = "(no output)"
	}
	return fmt.Sprintf("%s %s", r.Origin, safe), nil
}
