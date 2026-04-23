package main

import (
	"fmt"
	"os"
	"strings"

	"skyra-v04/src/primitives/entity"
	_ "skyra-v04/src/primitives/medium"
	"skyra-v04/src/primitives/world"
)

func main() {
	w, _ := world.World{}.Relate(entity.Relation{}).(world.World)

	if err := bootstrap(w); err != nil {
		fmt.Fprintln(os.Stderr, "bootstrap:", err)
		os.Exit(1)
	}

	// Kick off: skyra starts a thread with michael. start-thread creates the thread,
	// registers it, and routes the first continue-thread.
	rel, err := entity.Impress("michael", "", "skyra start-thread ~with skyra ~about conversation ~because bootstrap ~say hi | main")
	if err != nil {
		fmt.Fprintln(os.Stderr, "kickoff:", err)
		os.Exit(1)
	}
	node, ok := w.EntityMap[rel.ID]
	if !ok {
		fmt.Fprintln(os.Stderr, "start-thread not found")
		os.Exit(1)
	}
	node.Relate(rel)
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
		rel, err := entity.Impress("genome", "", line)
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
