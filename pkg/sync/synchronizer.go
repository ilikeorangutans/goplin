package sync

import (
	"log"
	"path/filepath"

	"github.com/ilikeorangutans/goplin/pkg/database"
	"github.com/pkg/errors"
)

func NewSynchronizer(syncDirPath string, database *database.Database) (*Synchronizer, error) {
	dir, err := OpenDir(syncDirPath)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create synchronizer")
	}
	lockDir, err := NewLock(filepath.Join(syncDirPath, ".lock"))
	if err != nil {
		return nil, errors.Wrap(err, "cannot create synchronizer")
	}

	return &Synchronizer{
		dir:      dir,
		lockDir:  lockDir,
		database: database,
	}, nil
}

type Synchronizer struct {
	dir      *Dir
	database *database.Database
	lockDir  *LockDir
}

func (s *Synchronizer) Sync() error {
	log.Println("starting sync")
	lock, err := s.lockDir.Acquire()
	if err != nil {
		return err
	}
	defer lock.Release()

	items, err := s.dir.Read()
	if err != nil {
		return err
	}
	for _, item := range items.Items {
		// Check if we have that item
		hasID, err := s.database.HasSyncItem(item.ID)
		if err != nil {
			return err
		}

		if hasID {
			log.Printf("update to existing id %s", item.ID)
		} else {
			log.Printf("new sync item %s", item.ID)
			tx := s.database.Begin()
			tx.MustExec("insert into sync_items (id) values ($1)", item.ID)
			// begin transaction
			// 1. insert the actual item in the database
			// 2. copy over resource if necessary
			// 3. insert into sync items
			// commit
			err = tx.Commit()
			if err != nil {
				return err
			}
		}

		// If we have it, check if it has changed
	}

	return nil
}
