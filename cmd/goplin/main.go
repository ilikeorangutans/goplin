package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/ilikeorangutans/goplin/pkg/change"
	"github.com/ilikeorangutans/goplin/pkg/cmdbar"
	"github.com/ilikeorangutans/goplin/pkg/database"
	"github.com/ilikeorangutans/goplin/pkg/model"
	"github.com/ilikeorangutans/goplin/pkg/sync"
	"github.com/rivo/tview"
)

type Command struct {
	Execute func(notebooks *model.NotebookService) error
}

func NotebookToTreeNode(notebook *model.Notebook) *tview.TreeNode {
	node := tview.NewTreeNode(notebook.Title)

	for _, child := range notebook.Notebooks {
		childNode := NotebookToTreeNode(child)
		node.AddChild(childNode)
	}
	node.SetReference(notebook.ID)
	return node
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	app := tview.NewApplication()
	treeView := tview.NewTreeView()

	notebooks := model.NewNotebookService()
	notebooks.Create("stuff", nil)
	notebooks.Create("more stuff", nil)
	notebooks.Create("some stuff", nil)
	notebooks.Create("no stuff", nil)
	notebooks.Create("only stuff", nil)

	commands := make(chan Command)
	go func(ctx context.Context, commands chan Command) {
		for {
			select {
			case <-ctx.Done():
				return
			case cmd := <-commands:
				err := cmd.Execute(notebooks)
				if err != nil {
					panic(err)
				}

				app.QueueUpdateDraw(func() {
					treeView.GetRoot().ClearChildren()
					for _, notebook := range notebooks.RootNotebooks() {
						node := NotebookToTreeNode(notebook)

						treeView.GetRoot().AddChild(node)
					}

					treeView.SetCurrentNode(treeView.GetRoot().GetChildren()[0])
				})
			}
		}
	}(ctx, commands)

	root := tview.NewFlex()
	root.SetDirection(tview.FlexRow)
	mainPanel := tview.NewFlex()
	logWindow := tview.NewTextView()

	notes := tview.NewList()
	rootNode := tview.NewTreeNode("<root>")
	treeView.SetRoot(rootNode)
	treeView.SetTopLevel(1)
	treeView.SetBorder(true).SetTitle("Notebooks")
	treeView.SetSelectedFunc(func(node *tview.TreeNode) {
		logWindow.SetText(logWindow.GetText(false) + fmt.Sprintf("Selected tree node %s: %s\n", node.GetText(), node.GetReference()))
		logWindow.ScrollToEnd()
	})
	treeView.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyTab {
			app.SetFocus(notes)
		}
	})
	notes.SetDoneFunc(func() {
		app.SetFocus(treeView)
	})

	mainPanel.AddItem(treeView, 0, 1, false)

	notes.SetBorder(true).SetTitle("Notes")

	mainPanel.AddItem(notes, 0, 1, false)

	noteDetails := tview.NewFlex()

	mainPanel.AddItem(noteDetails, 0, 2, false)

	logWindow.SetBorder(true).SetTitle("logs")
	logWindow.SetScrollable(true)

	root.AddItem(mainPanel, 0, 3, false)
	root.AddItem(logWindow, 10, 1, false)

	cmdBar := cmdbar.New("begin entering commands by pressing \":\" or begin search with \"/\"")
	cmdBar.AddCmd(":q", "quit the app", func(_ string) error {
		cancel()
		return nil
	})
	cmdBar.AddCmd(":mkbook", "make a top level notebook, prefix with // to place under current notebook", func(name string) error {
		var err error
		var parent *model.Notebook
		if strings.HasPrefix(name, "//") {
			parent, err = notebooks.NotebookByID(treeView.GetCurrentNode().GetReference().(string))
			if err != nil {
				logWindow.SetText(logWindow.GetText(false) + "couldn't find parent: " + err.Error())
				logWindow.ScrollToEnd()
				return err
			}
			logWindow.SetText(logWindow.GetText(false) + "adding notebook with parent " + parent.Title)
			logWindow.ScrollToEnd()
			name = name[2:]
		}

		logWindow.SetText(logWindow.GetText(false) + "adding notebook " + name)
		logWindow.ScrollToEnd()
		commands <- Command{Execute: func(notebooks *model.NotebookService) error {
			change.AddNotebook(name, parent).Apply(notebooks)
			return nil
		}}
		cmdBar.SetText("")
		app.SetFocus(treeView)

		return nil
	})
	cmdBar.AddCmd(":rmbook", "remove currently selected notebook", func(_ string) error {
		id := treeView.GetCurrentNode().GetReference().(string)
		logWindow.SetText(logWindow.GetText(false) + "deleting notebook " + id)
		logWindow.ScrollToEnd()
		commands <- Command{Execute: func(notebooks *model.NotebookService) error {
			change.DeleteNotebook(id).Apply(notebooks)
			return nil
		}}
		cmdBar.SetText("")
		app.SetFocus(treeView)
		return nil
	})
	cmdBar.AddCmd(":cmds", "shows all known commands and their usage", func(_ string) error {
		var sbuilder strings.Builder
		fmt.Fprint(&sbuilder, "Known commands:\n")
		for _, c := range cmdBar.SummarizeCommands() {
			fmt.Fprintf(&sbuilder, "%-10s  %s\n", c.Verb, c.Summary)
		}
		logWindow.SetText(sbuilder.String())
		logWindow.ScrollToEnd()
		return nil
	})

	root.AddItem(cmdBar, 1, 1, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if app.GetFocus() != cmdBar && event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case ':':
				cmdBar.SetText(":")
				app.SetFocus(cmdBar)
			case '/':
				cmdBar.SetText("/")
				app.SetFocus(cmdBar)
			}
		} else if event.Key() == tcell.KeyEscape && app.GetFocus() == cmdBar {
			app.SetFocus(treeView)
		}
		return event
	})
	go func() {
		for {
			select {
			case <-ctx.Done():
				app.Stop()

			}
		}
	}()

	var sbuilder strings.Builder
	fmt.Fprint(&sbuilder, "Known commands:\n")
	for _, c := range cmdBar.SummarizeCommands() {
		fmt.Fprintf(&sbuilder, "%-10s  %s\n", c.Verb, c.Summary)
	}
	logWindow.SetText(sbuilder.String())
	logWindow.ScrollToEnd()

	app.SetRoot(root, true)
	app.SetFocus(treeView)
	// build initial ui
	for _, notebook := range notebooks.RootNotebooks() {
		node := NotebookToTreeNode(notebook)

		treeView.GetRoot().AddChild(node)
	}

	treeView.SetCurrentNode(treeView.GetRoot().GetChildren()[0])

	err := app.Run()
	if err != nil {
		panic(err)
	}
}

func foo() {
	db, err := database.OpenDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Migrate()
	if err != nil {
		log.Fatal(err)
	}
	syncDir, err := sync.OpenDir(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	items, err := syncDir.Read()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("got %d items from sync dir", items.Len())

}
