package keychain

import (
	"os/exec"
	"os/user"
	"strings"
)

// Get reads a secret from the macOS Keychain. Falls back to empty string on error.
// Store a secret once with:
//
//	security add-generic-password -a "$USER" -s <service> -w
func Get(service string) string {
	u, err := user.Current()
	if err != nil {
		return ""
	}
	out, err := exec.Command("security", "find-generic-password", "-a", u.Username, "-s", service, "-w").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
