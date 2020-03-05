package change

import (
	"github.com/ilikeorangutans/goplin/pkg/model"
	"github.com/pkg/errors"
)

type Change interface {
	ItemID() string
	Apply(*model.NotebookService) error // TODO maybe return a boolean if something changed?
}

type base struct {
	itemID string
}

func (b base) ItemID() string {
	return b.itemID
}

type addNotebook struct {
	base
	name     string
	parentID string
}

func (a *addNotebook) Apply(service *model.NotebookService) error {
	var err error
	var parent *model.Notebook
	if a.parentID != "" {
		parent, err = service.NotebookByID(a.parentID)
		if err != nil {
			return errors.Wrap(err, "could not find parent id")
		}
	}

	_, err = service.Create(a.name, parent)
	return err
}

func AddNotebook(title string, parent *model.Notebook) Change {
	var parentID string
	if parent != nil {
		parentID = parent.ID
	}
	return &addNotebook{
		name:     title,
		parentID: parentID,
	}
}

type deleteNotebook struct {
	base
}

func (d *deleteNotebook) Apply(service *model.NotebookService) error {
	return service.DeleteNotebook(d.itemID)
}

func DeleteNotebook(id string) Change {
	return &deleteNotebook{
		base: base{
			itemID: id,
		},
	}
}
