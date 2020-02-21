package main

import (
	"log"
	"os"

	"github.com/ilikeorangutans/goplin/pkg/sync"
)

func main() {
	syncDir, err := sync.OpenDir(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	items, err := syncDir.Read()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Got %d items", len(items))
	for _, item := range items {
		if item.Type == sync.TypeNote {
			log.Printf("Note ")
		}
		if item.Type == sync.TypeFolder {
			log.Printf("Notebook %s", item.Body)
		}
	}
}
