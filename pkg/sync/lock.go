package sync

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
)

func NewLock(path string) (*LockDir, error) {
	if fileInfo, err := os.Stat(path); err != nil {
		return nil, err
	} else if !fileInfo.IsDir() {
		return nil, fmt.Errorf("%q is not a directory", path)
	}

	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return &LockDir{
		dir:       dir,
		path:      path,
		generator: func() string { return fmt.Sprintf("goplin_%d", time.Now().UnixNano()/int64(time.Millisecond)) },
	}, nil
}

type LockNameGenerator func() string

// LockDir represents the lock directory and operations to acquire locks in it.
type LockDir struct {
	dir       *os.File
	path      string
	generator LockNameGenerator
}

func (l *LockDir) Acquire() (*Lock, error) {
	locked, err := l.IsLocked()
	if err != nil {
		return nil, err
	}
	if locked {
		return nil, fmt.Errorf("directory is locked")
	}

	lockPath := filepath.Join(l.path, l.generator())
	f, err := os.Create(lockPath)
	if err != nil {
		return nil, err
	}
	err = f.Close()
	if err != nil {
		return nil, err
	}
	return &Lock{
		path: lockPath,
	}, nil
}

func (l *LockDir) IsLocked() (bool, error) {
	entries, err := l.dir.Readdir(0)
	if err != nil {
		return false, errors.Wrap(err, "could not read lock dir")
	}
	return len(entries) > 0, nil
}

type Lock struct {
	path string
}

func (l *Lock) Release() error {
	log.Printf("releasing lock %s", l.path)
	return os.Remove(l.path)
}
