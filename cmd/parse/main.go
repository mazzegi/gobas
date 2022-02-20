package main

import (
	"fmt"
	"time"

	"github.com/mazzegi/gobas"
)

func main() {
	parser := gobas.NewParser()
	//_, err := parser.ParseFile("../../samples/001_aceyducey.bas")
	t0 := time.Now()
	_, err := parser.ParseFile("../../samples/002_amazing.bas")
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	fmt.Printf("parsed in %s\n", time.Since(t0))

	// stmts, err := gobas.ParseFile("../../samples/001_aceyducey.bas")
	// if err != nil {
	// 	panic(err)
	// }
	// for _, stmt := range stmts {
	// 	js, _ := json.Marshal(stmt)
	// 	fmt.Printf("%T: %s\n", stmt, string(js))
	// }
}
