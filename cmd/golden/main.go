package main

import (
	"fmt"
	"golden/internal/golden"
	"log"
)

const defaultBindAddr = "0.0.0.0:8091"

func main() {
	fmt.Println("GoldenDB is booting")

	g := golden.Golden{
		ObjLimit:   65535,             // 64K objects count limit
		MemLimit:   128 * 1024 * 1024, // 128 MB per object size max limit
		MaxWorkers: 200,               // default limit
	}

	log.Printf("Binding with TCP addr %s and waiting for connections", defaultBindAddr)
	if err := g.Bind(defaultBindAddr); err != nil {
		log.Fatal(err.Error())
	}
}
