package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/ilikeorangutans/goplin/pkg/change"
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

	inputBox := tview.NewInputField()
	inputBox.SetPlaceholder("begin entering commands by pressing \":\" or begin search with \"/\"")
	inputBox.SetChangedFunc(func(input string) {
		s := strings.TrimSpace(input)
		logWindow.SetText(logWindow.GetText(false) + s)
		logWindow.ScrollToEnd()
	})
	inputBox.SetDoneFunc(func(key tcell.Key) {
		if key != tcell.KeyEnter {
			inputBox.SetText("")
			return
		}
		input := strings.TrimSpace(inputBox.GetText())

		if input == ":q" || input == ":quit" {
			cancel()
		} else if strings.HasPrefix(input, ":mkbook ") {
			name := strings.SplitN(input, " ", 2)[1]

			var err error
			var parent *model.Notebook
			if strings.HasPrefix(name, "//") {
				parent, err = notebooks.NotebookByID(treeView.GetCurrentNode().GetReference().(string))
				if err != nil {
					logWindow.SetText(logWindow.GetText(false) + "couldn't find parent: " + err.Error())
					logWindow.ScrollToEnd()
					return
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
			inputBox.SetText("")
			app.SetFocus(treeView)
		} else if strings.HasPrefix(input, ":rmbook") {
			id := treeView.GetCurrentNode().GetReference().(string)
			logWindow.SetText(logWindow.GetText(false) + "deleting notebook " + id)
			logWindow.ScrollToEnd()
			commands <- Command{Execute: func(notebooks *model.NotebookService) error {
				change.DeleteNotebook(id).Apply(notebooks)
				return nil
			}}
			inputBox.SetText("")
			app.SetFocus(treeView)

		}
	})
	root.AddItem(inputBox, 1, 1, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if app.GetFocus() != inputBox && event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case ':':
				inputBox.SetText(":")
				app.SetFocus(inputBox)
			case '/':
				inputBox.SetText("/")
				app.SetFocus(inputBox)
			}
		} else if event.Key() == tcell.KeyEscape && app.GetFocus() == inputBox {
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

	logWindow.SetText(`
known commands:
q, quit          quit
mkbook title     creates a new notebook at the root level
mkbook //title   creates a notebook under the currently selected notebook
rmbook           deletes the currently selected notebook
	`)
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
