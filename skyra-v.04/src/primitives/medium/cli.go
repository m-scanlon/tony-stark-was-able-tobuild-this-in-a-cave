package medium

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"skyra-v04/src/primitives/entity"
)

var stdinScanner = bufio.NewScanner(os.Stdin)

func init() {
	Register("cli", cli)
}

func cli(present string, r entity.Relation) (string, error) {
	fmt.Println("\n---")
	fmt.Print(present)
	fmt.Println("\n---")
	fmt.Print("> ")
	if !stdinScanner.Scan() {
		return "", stdinScanner.Err()
	}
	input := strings.TrimSpace(stdinScanner.Text())
	if input == "" {
		return "", nil
	}
	// If the user typed a valid protocol line, pass it through. Otherwise wrap it
	// as a continue-thread reply to whoever is in the current exchange with us.
	if _, err := entity.Impress("", "", input); err == nil {
		return input, nil
	}
	return fmt.Sprintf("skyra continue-thread ~with %s ~say %s | user", r.Origin, input), nil
}
