package sync

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/pkg/errors"
)

var ItemRegex = regexp.MustCompile("[a-z0-9]{32}.md")

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

	lockDir, err := NewLock(filepath.Join(path, ".lock"))
	return &SyncDir{
		dir:     dir,
		path:    path,
		lockDir: lockDir,
	}, nil
}

type SyncDir struct {
	dir     *os.File
	path    string
	lockDir *LockDir
}

func (sd *SyncDir) Read() (*Items, error) {
	lock, err := sd.lockDir.Acquire()
	if err != nil {
		return nil, errors.Wrap(err, "could not acquire lock")
	}
	defer lock.Release()
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
