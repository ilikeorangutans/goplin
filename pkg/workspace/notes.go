package workspace

import (
	"fmt"

	"github.com/ilikeorangutans/goplin/pkg/model"
)

func NewNotes() *Notes {
	return &Notes{
		byNotebook: make(map[string][]*model.Note),
		byID:       make(map[string]*model.Note),
		Events:     make(chan Event),
	}
}

type Notes struct {
	notes      []*model.Note
	byNotebook map[string][]*model.Note
	byID       map[string]*model.Note
	Events     chan Event
}

func (n *Notes) ByID(id string) (*model.Note, error) {
	note, ok := n.byID[id]
	if !ok {
		return nil, ErrNoSuchItem
	}
	return note, nil
}

func (n *Notes) Create(title string, notebook *model.Notebook) (*model.Note, error) {
	if notebook == nil {
		return nil, fmt.Errorf("cannot create note without notebook")
	}
	item := model.NewItem().WithParent(notebook.Item)

	note := &model.Note{
		Item:  item,
		Title: title,
	}

	return n.Save(note)
}

func (n *Notes) ByNotebook(notebook *model.Notebook) ([]*model.Note, error) {
	return n.byNotebook[notebook.ID], nil
}

func (n *Notes) Save(note *model.Note) (*model.Note, error) {
	n.notes = append(n.notes, note)
	n.byNotebook[note.ParentID] = append(n.byNotebook[note.ParentID], note)
	n.byID[note.ID] = note

	n.Events <- Event{ID: note.ID}

	return note, nil
}
