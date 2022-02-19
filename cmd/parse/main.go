package main

import (
	"encoding/json"
	"fmt"

	"github.com/mazzegi/gobas"
)

func main() {
	stmts, err := gobas.ParseFile("../../samples/001_aceyducey.bas")
	if err != nil {
		panic(err)
	}
	for _, stmt := range stmts {
		js, _ := json.Marshal(stmt)
		fmt.Printf("%T: %s\n", stmt, string(js))
	}
}
