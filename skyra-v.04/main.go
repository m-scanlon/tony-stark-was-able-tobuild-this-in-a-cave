package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"skyra-v04/src/primitives/entity"
	_ "skyra-v04/src/primitives/medium"
	"skyra-v04/src/primitives/world"
)

func main() {
	loadEnv("../.env")
	_ = os.RemoveAll("logs")

	w := world.New()

	if err := bootstrap(w); err != nil {
		fmt.Fprintln(os.Stderr, "bootstrap:", err)
		os.Exit(1)
	}

	rel, err := entity.Impress("michael", "", "skyra hi ~about conversation ~because bootstrap")
	if err != nil {
		fmt.Fprintln(os.Stderr, "kickoff:", err)
		os.Exit(1)
	}
	w.Relate(rel)
}

func loadEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		if os.Getenv(k) == "" {
			os.Setenv(k, v)
		}
	}
}

func bootstrap(w *world.World) error {
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
		w.Relate(rel)
	}
	return nil
}
