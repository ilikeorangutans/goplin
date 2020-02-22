package sync

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/ilikeorangutans/goplin/pkg/model"
)

type Tag struct {
}

func NewNote(id string, items *Items) (*model.Note, error) {
	note, err := items.FindByID(id)
	if err != nil {
		return nil, err
	}
	if note.Type != TypeNote {
		return nil, fmt.Errorf("tried to open item id %q as a Note but it is of type %s (%d)", id, note.Type, note.Type)
	}
	notebook, err := NewNotebook(note.ParentID, items)
	if err != nil {
		return nil, err
	}
	return &model.Note{
		Item: model.Item{
			ID:       note.ID,
			ParentID: note.ParentID,
		},
		Title:    titleFromBody(note.Body),
		Body:     note.Body,
		Notebook: notebook,
	}, nil
}

func titleFromBody(body string) string {
	s := bufio.NewScanner(strings.NewReader(body))
	s.Scan()
	return strings.TrimSpace(s.Text())
}

func NotebooksFromSyncDir(items *Items) (*model.Notebooks, error) {
	var notebooks []*model.Notebook
	byID := make(map[string]*model.Notebook)
	byParentID := make(map[string][]*model.Notebook)
	for _, item := range items.Items {
		if item.Type != TypeFolder {
			continue
		}

		notebook, err := NewNotebook(item.ID, items)
		if err != nil {
			return nil, err
		}

		byID[item.ID] = notebook
		if item.ParentID != "" {
			byParentID[notebook.ParentID] = append(byParentID[notebook.ParentID], notebook)
		}

		notebooks = append(notebooks, notebook)
	}

	queue := make([]*model.Notebook, len(notebooks))
	copy(queue, notebooks)

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:len(queue)]

		children := byParentID[cur.ID]
		queue = append(queue, children...)

		cur.Notebooks = children
	}

	return &model.Notebooks{
		Notebooks: notebooks,
		ByID:      byID,
	}, nil
}

func NewNotebook(id string, items *Items) (*model.Notebook, error) {
	if id == "" {
		return nil, fmt.Errorf("trying to create notebook with empty id")
	}
	note, err := items.FindByID(id)
	if err != nil {
		return nil, err
	}
	if note.Type != TypeFolder {
		return nil, fmt.Errorf("tried to open item id %q as a Folder but it is of type %s (%d)", id, note.Type, note.Type)
	}

	return &model.Notebook{
		Item: model.Item{
			ID:       note.ID,
			ParentID: note.ParentID,
		},
		Title: titleFromBody(note.Body),
		Order: note.Order,
	}, nil
}
