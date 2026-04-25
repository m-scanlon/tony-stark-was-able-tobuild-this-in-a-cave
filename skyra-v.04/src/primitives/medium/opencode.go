package medium

import (
	"fmt"
	"os/exec"
	"strings"

	"skyra-v04/src/primitives/entity"
)

func init() {
	Register("opencode", opencodeMedium)
}

func opencodeMedium(_ string, r entity.Relation) (string, error) {
	msg := strings.TrimSpace(r.Impulse)
	if msg == "" {
		return "", nil
	}
	out, err := exec.Command("/opt/homebrew/bin/opencode", "run", "--format", "json", msg).CombinedOutput()
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
