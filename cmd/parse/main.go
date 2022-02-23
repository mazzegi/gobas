package main

import (
	"fmt"
	"time"

	"github.com/mazzegi/gobas"
)

func main() {
	parser := gobas.NewParser()
	t0 := time.Now()
	stmts, err := parser.ParseFile("../../samples/001_aceyducey.bas")
	//_, err := parser.ParseFile("../../samples/002_amazing.bas")
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	fmt.Printf("parsed in %s\n", time.Since(t0))
	_ = stmts
	// for _, stmt := range stmts {
	// 	js, _ := json.Marshal(stmt)
	// 	fmt.Printf("%T: %s\n", stmt, string(js))
	// }
}
