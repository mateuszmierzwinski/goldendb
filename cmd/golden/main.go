package main

import (
	"fmt"
	"golden/internal/golden"
	"log"
)

func main() {
	fmt.Println("GoldenDB is booting")

	g := golden.Golden{
		ObjLimit:   65535,             // 64K objects count limit
		MemLimit:   128 * 1024 * 1024, // 128 MB per object size max limit
		MaxWorkers: 200,               // default limit
	}

	log.Println("Binding with TCP addr 0.0.0.0:8091 and waiting for connections")
	if err := g.Bind("0.0.0.0:8091"); err != nil {
		log.Fatal(err.Error())
	}
}
