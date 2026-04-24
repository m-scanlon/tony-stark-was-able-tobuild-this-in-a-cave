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
	if _, err := entity.Impress("", "", input); err == nil {
		return input, nil
	}
	return fmt.Sprintf("%s %s", r.Origin, input), nil
}
