package model

import "fmt"

type Note struct {
	Item
	Title    string
	Body     string
	Notebook *Notebook
}

func NewNotebookService() *NotebookService {
	return &NotebookService{
		notebooksByID: make(map[string]*Notebook),
	}
}

type NotebookService struct {
	notebooks     []*Notebook
	notebooksByID map[string]*Notebook
}

var ErrNoSuchItem = fmt.Errorf("no such item")

func (n *NotebookService) NotebookByID(id string) (*Notebook, error) {
	notebook, ok := n.notebooksByID[id]
	if !ok {
		return nil, ErrNoSuchItem
	}
	return notebook, nil
}

func (n *NotebookService) DeleteNotebook(id string) error {
	_, err := n.NotebookByID(id)
	if err != nil {
		return err
	}

	i := 0
	for _, x := range n.notebooks {
		if x.ID != id {
			n.notebooks[i] = x
			i++
		}

	}
	n.notebooks = n.notebooks[:i]
	// TODO need ot remove all other references in the child nodes too

	delete(n.notebooksByID, id)
	return nil
}

func (n *NotebookService) Create(name string, parent *Notebook) (*Notebook, error) {
	item := NewItem()
	if parent != nil {
		item = item.WithParent(parent.Item)
	}

	notebook := &Notebook{
		Item:   item,
		Title:  name,
		Parent: parent,
	}

	n.notebooks = append(n.notebooks, notebook)
	n.notebooksByID[item.ID] = notebook
	if parent != nil {
		parent.Notebooks = append(parent.Notebooks, notebook)
	}
	return notebook, nil
}

func (n *NotebookService) TopLevel() []*Notebook {
	var result []*Notebook

	for _, notebook := range n.notebooks {
		if notebook.HasParent() {
			continue
		}
		result = append(result, notebook)
	}

	return result
}
