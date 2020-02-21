package main

import (
	"log"

	"github.com/ilikeorangutans/goplin/pkg/sync"
)

func main() {
	syncDir, err := sync.OpenDir("/home/jakob/notes")
	if err != nil {
		log.Fatal(err)
	}
	err = syncDir.Read()
	if err != nil {
		log.Fatal(err)
	}
}
