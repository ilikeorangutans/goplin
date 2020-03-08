package tui

import "github.com/rivo/tview"

type UIUpdater interface {
	QueueUpdateDraw(func())
}

type Command interface {
	Execute() error
}

type FooCommand struct {
	UIUpdater UIUpdater
	Tree      *tview.TreeView
}

func (f *FooCommand) Execute() error {
	f.UIUpdater.QueueUpdateDraw(func() {
		node := tview.NewTreeNode("fooooo")
		f.Tree.GetRoot().AddChild(node)
		f.Tree.SetCurrentNode(node)
	})
	return nil
}
