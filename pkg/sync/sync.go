package sync

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func OpenDir(path string) (*SyncDir, error) {
	if fileInfo, err := os.Stat(path); err != nil {
		return nil, err
	} else if !fileInfo.IsDir() {
		return nil, fmt.Errorf("%q is not a directory", path)
	}

	dir, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	return &SyncDir{
		dir:  dir,
		path: path,
	}, nil
}

type SyncDir struct {
	dir  *os.File
	path string
}

var ItemRegex = regexp.MustCompile("[a-z0-9]{32}.md")

func (sd *SyncDir) Read() error {
	// TODO this should probably check the lock directory
	log.Printf("reading %s", sd.dir.Name())
	entries, err := sd.dir.Readdir(0)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if ItemRegex.MatchString(entry.Name()) {
			log.Printf("Item %s", entry.Name())
			item, err := ReadItem(filepath.Join(sd.path, entry.Name()))
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("%v", item)
		}
	}
	return nil
}
