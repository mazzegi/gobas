package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mazzegi/gobas"
)

func main() {
	dir, _ := filepath.Abs("../../100games/")
	fis, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	var resl []string
	for _, fi := range fis {
		if filepath.Ext(fi.Name()) != ".bas" {
			continue
		}

		path := filepath.Join(dir, fi.Name())
		parser := gobas.NewParser()
		t0 := time.Now()
		_, err := parser.ParseFile(path)
		if err != nil {
			resl = append(resl, fmt.Sprintf("FAILED: %q: %v", path, err))
		} else {
			resl = append(resl, fmt.Sprintf("OK: %q => %s", path, time.Since(t0)))
		}
	}

	fmt.Println("*** result ***")
	for _, s := range resl {
		fmt.Println(s)
	}
}
