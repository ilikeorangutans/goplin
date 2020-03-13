package workspace

import (
	"fmt"

	"github.com/ilikeorangutans/goplin/pkg/model"
)

var ErrNoSuchItem = fmt.Errorf("no such item")

func NewNotebooks() *Notebooks {
	return &Notebooks{
		byID: make(map[string]*model.Notebook),
	}
}

// TODO add a callback mechanism so I can easily listen to changes in notebooks
type Notebooks struct {
	notebooks []*model.Notebook
	byID      map[string]*model.Notebook
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
