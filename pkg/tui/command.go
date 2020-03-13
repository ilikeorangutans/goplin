package tui

import (
	"github.com/ilikeorangutans/goplin/pkg/model"
	"github.com/ilikeorangutans/goplin/pkg/workspace"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
)

type UIUpdater interface {
	QueueUpdateDraw(func())
}

type Command interface {
	Execute() error
}

type CreateNotebookCommand struct {
	UIUpdater
	Tree      *tview.TreeView
	Workspace *workspace.Workspace
	Name      string
	Parent    *model.Notebook
}

func (c *CreateNotebookCommand) Execute() error {
	notebook, err := c.Workspace.Notebooks().Create(c.Name, c.Parent)
	if err != nil {
		return errors.Wrap(err, "could not create notebook")
	}
	c.QueueUpdateDraw(func() {
		node := tview.NewTreeNode(notebook.Title)
		node.SetReference(notebook)
		c.Tree.GetRoot().AddChild(node)
		c.Tree.SetCurrentNode(node)
	})
	return nil
}

type CreateNoteCommand struct {
	UIUpdater
	Notes     *tview.List
	Workspace *workspace.Workspace
	Name      string
	Parent    *model.Notebook
}

func (c *CreateNoteCommand) Execute() error {
	note, err := c.Workspace.Notes().Create(c.Name, c.Parent)
	if err != nil {
		return errors.Wrap(err, "could not create note")
	}
	c.QueueUpdateDraw(func() {
		c.Notes.AddItem(note.Title, "", -1, func() {})
	})

	return nil
}
