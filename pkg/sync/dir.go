package sync

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var ItemRegex = regexp.MustCompile("[a-z0-9]{32}.md")

func OpenDir(path string) (*Dir, error) {
	if fileInfo, err := os.Stat(path); err != nil {
		return nil, err
	} else if !fileInfo.IsDir() {
		return nil, fmt.Errorf("%q is not a directory", path)
	}

	dir, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	return &Dir{
		dir:  dir,
		path: path,
	}, nil
}

type Dir struct {
	dir  *os.File
	path string
}

func (sd *Dir) Read() (*Items, error) {
	log.Printf("reading %s", sd.dir.Name())
	entries, err := sd.dir.Readdir(0)
	if err != nil {
		return nil, err
	}

	items := make([]*Item, 0, len(entries))

	for _, entry := range entries {
		if ItemRegex.MatchString(entry.Name()) {
			item, err := ReadItem(filepath.Join(sd.path, entry.Name()))
			if err != nil {
				return nil, err
			}
			items = append(items, item)
		}
	}
	return NewItems(items), nil
}
