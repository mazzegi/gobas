package main

import (
	"fmt"

	"github.com/mazzegi/gobas"
)

func main() {
	parser := gobas.NewParser()
	//state, err := parser.ParseFile("../../samples/001_aceyducey.bas")
	//state, err := parser.ParseFile("../../samples/002_amazing_debug.bas")
	//state, err := parser.ParseFile("../../samples/002_amazing.bas")
	state, err := parser.ParseFile("../../100games/003_animal.bas")
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	state.Run()
}
