package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"skyra-v04/src/primitives/entity"
	"skyra-v04/src/primitives/world"
)

func main() {
	w, _ := world.World{}.Relate(entity.Relation{}).(world.World)

	if err := bootstrap(w); err != nil {
		fmt.Fprintln(os.Stderr, "bootstrap:", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		raw := strings.TrimSpace(scanner.Text())
		if raw == "" {
			fmt.Print("> ")
			continue
		}

		rel, err := entity.Parse("michael", "main", "skyra continue-thread ~with skyra ~say "+raw+" | main")
		if err != nil {
			fmt.Println("error:", err)
			fmt.Print("> ")
			continue
		}

		node, ok := w.EntityMap[rel.ID]
		if !ok {
			fmt.Println("error: unknown target:", rel.ID)
			fmt.Print("> ")
			continue
		}

		node.Relate(rel)
		fmt.Print("> ")
	}
}

func bootstrap(w world.World) error {
	data, err := os.ReadFile("genome.skyra")
	if err != nil {
		return err
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		rel, err := entity.Parse("genome", "", line)
		if err != nil {
			return fmt.Errorf("genome: %w", err)
		}
		node, ok := w.EntityMap[rel.ID]
		if !ok {
			return fmt.Errorf("genome: unknown target %q", rel.ID)
		}
		node.Relate(rel)
	}
	return nil
}
