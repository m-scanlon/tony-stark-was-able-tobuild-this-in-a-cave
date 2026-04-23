package medium

import (
	"fmt"
	"os/exec"
	"strings"

	"skyra-v04/src/primitives/entity"
	"skyra-v04/src/primitives/meaning"
)

func init() {
	Register("shell", shell)
}

func shell(_ string, r entity.Relation) (string, error) {
	cmd, err := meaning.Extract(r.Impulse, "~say", "shell", "|")
	if err != nil {
		return "", nil
	}
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	response := string(out)
	if err != nil {
		response = "error: " + err.Error() + "\n" + response
	}
	// Sanitize so the output survives Impress: escape pipes, collapse newlines to a visible marker.
	safe := strings.ReplaceAll(response, "|", "│")
	safe = strings.ReplaceAll(safe, "\n", " ↵ ")
	safe = strings.TrimSpace(safe)
	if safe == "" {
		safe = "(no output)"
	}
	return fmt.Sprintf("skyra continue-thread ~with %s ~say %s | shell-out", r.Origin, safe), nil
}
