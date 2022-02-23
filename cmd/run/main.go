package main

import (
	"fmt"

	"github.com/mazzegi/gobas"
)

func main() {
	parser := gobas.NewParser()
	state, err := parser.ParseFile("../../samples/001_aceyducey.bas")
	//_, err := parser.ParseFile("../../samples/002_amazing.bas")
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	state.Run()
}
