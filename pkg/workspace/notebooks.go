package workspace

import (
	"fmt"

	"github.com/ilikeorangutans/goplin/pkg/model"
)

var ErrNoSuchItem = fmt.Errorf("no such item")

type EventType int

func NewNotebooks() *Notebooks {
	return &Notebooks{
		byID:   make(map[string]*model.Notebook),
		Events: make(chan Event),
	}
}

type Notebooks struct {
	notebooks []*model.Notebook
	byID      map[string]*model.Notebook
	Events    chan Event
}

func (n *Notebooks) ByID(id string) (*model.Notebook, error) {
	notebook, ok := n.byID[id]
	if !ok {
		return nil, ErrNoSuchItem
	}
	return notebook, nil
}

// TopLevel returns all notebooks that have no parent.
func (n *Notebooks) TopLevel() ([]*model.Notebook, error) {
	var result []*model.Notebook

	for _, notebook := range n.notebooks {
		if notebook.HasParent() {
			continue
		}
		result = append(result, notebook)
	}

	return result, nil
}

func (n *Notebooks) Delete(id string) error {
	// TODO implement me
	return nil
}

func (n *Notebooks) Save(notebook *model.Notebook) (*model.Notebook, error) {
	n.notebooks = append(n.notebooks, notebook)
	n.byID[notebook.ID] = notebook
	// TODO we should differentiate between update and create
	n.Events <- Event{ID: notebook.ID}
	return notebook, nil
}

// Create initializes a new notebook with the given name and the optional parent.
// The method returns the new notebook instance.
func (n *Notebooks) Create(name string, parent *model.Notebook) (*model.Notebook, error) {
	item := model.NewItem()
	if parent != nil {
		item = item.WithParent(parent.Item)
	}

	notebook := &model.Notebook{
		Item:   item,
		Title:  name,
		Parent: parent,
	}

	// TODO parent needs to be updated here?

	return n.Save(notebook)
}
