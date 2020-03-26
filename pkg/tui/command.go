package tui

import (
	"github.com/ilikeorangutans/goplin/pkg/model"
	"github.com/ilikeorangutans/goplin/pkg/workspace"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
)

type Command interface {
	Execute() error
}

type CreateNotebookCommand struct {
	Tree      *tview.TreeView
	Workspace *workspace.Workspace
	Name      string
	Parent    *model.Notebook
}

func (c *CreateNotebookCommand) Execute() error {
	_, err := c.Workspace.Notebooks().Create(c.Name, c.Parent)
	if err != nil {
		return errors.Wrap(err, "could not create notebook")
	}
	return nil
}

type CreateNoteCommand struct {
	Notes     *tview.List
	Workspace *workspace.Workspace
	Name      string
	Parent    *model.Notebook
}

func (c *CreateNoteCommand) Execute() error {
	_, err := c.Workspace.Notes().Create(c.Name, c.Parent)
	if err != nil {
		return errors.Wrap(err, "could not create note")
	}

	return nil
}
