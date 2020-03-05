package model

import (
	"strings"

	"github.com/google/uuid"
)

func NewItem() Item {
	id := strings.ReplaceAll(uuid.New().String(), "-", "")
	return Item{
		ID: id,
	}
}

type Item struct {
	ID       string
	ParentID string
	Source   string
}

func (i Item) WithParent(parent Item) Item {
	return Item{
		ID:       i.ID,
		ParentID: parent.ID,
	}
}

type ItemType int
