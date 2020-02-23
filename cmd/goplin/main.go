package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ilikeorangutans/goplin/pkg/model"
	"github.com/ilikeorangutans/goplin/pkg/sync"
)

func main() {
	db, err := model.OpenDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Migrate()
	if err != nil {
		log.Fatal(err)
	}
	syncDir, err := sync.OpenDir(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	items, err := syncDir.Read()
	if err != nil {
		log.Fatal(err)
	}

	notebooks, err := sync.NotebooksFromSyncDir(items)
	if err != nil {
		log.Fatal(err)
	}

	for _, notebook := range notebooks.Roots() {
		printNotebook(notebook, 0)
	}
}

func printNotebook(notebook *model.Notebook, level int) {
	indent := strings.Repeat("  ", level)
	fmt.Printf("%s%s\n", indent, notebook.Title)
	for _, child := range notebook.Notebooks {
		printNotebook(child, level+1)
	}
}
