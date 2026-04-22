package main

import (
	"bufio"
	"fmt"
	"os"

	"skyra-v04/src/primitives/logos"
	"skyra-v04/src/primitives/world"
)

func main() {
	w, _ := world.World{}.Relate(logos.Relation{}).(world.World)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		raw := scanner.Text()
		rel, err := logos.Parse("stdin", "", raw)
		if err != nil {
			fmt.Println("error:", err)
			fmt.Print("> ")
			continue
		}

		node, ok := w.LogosMap[rel.ID]
		if !ok {
			fmt.Println("error: unknown target:", rel.ID)
			fmt.Print("> ")
			continue
		}

		result := node.Relate(rel)
		fmt.Println(result.Name(), result.ID())
		fmt.Print("> ")
	}
}
