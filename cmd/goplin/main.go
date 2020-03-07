package main

import (
	"context"
	"log"
	"os"

	"github.com/ilikeorangutans/goplin/pkg/database"
	"github.com/ilikeorangutans/goplin/pkg/model"
	"github.com/ilikeorangutans/goplin/pkg/sync"
	"github.com/ilikeorangutans/goplin/pkg/tui"
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

	// TOOD create a workspace type
	notebooks := model.NewNotebookService()
	notebooks.Create("stuff", nil)
	notebooks.Create("more stuff", nil)
	notebooks.Create("some stuff", nil)
	notebooks.Create("no stuff", nil)
	notebooks.Create("only stuff", nil)

	ctx := context.Background()
	//  TODO inject workspace/notebooks into main
	mainUI := tui.NewMain()
	if err := mainUI.Run(ctx); err != nil {
		panic(err)
	}
}

func blarg() {
	//
	// 	commands := make(chan Command)
	// 	go func(ctx context.Context, commands chan Command) {
	// 		for {
	// 			select {
	// 			case <-ctx.Done():
	// 				return
	// 			case cmd := <-commands:
	// 				err := cmd.Execute(notebooks)
	// 				if err != nil {
	// 					panic(err)
	// 				}
	//
	// 				app.QueueUpdateDraw(func() {
	// 					treeView.GetRoot().ClearChildren()
	// 					for _, notebook := range notebooks.RootNotebooks() {
	// 						node := NotebookToTreeNode(notebook)
	//
	// 						treeView.GetRoot().AddChild(node)
	// 					}
	//
	// 					treeView.SetCurrentNode(treeView.GetRoot().GetChildren()[0])
	// 				})
	// 			}
	// 		}
	// 	}(ctx, commands)
	//
	// 	treeView.SetSelectedFunc(func(node *tview.TreeNode) {
	// 		logWindow.SetText(logWindow.GetText(false) + fmt.Sprintf("Selected tree node %s: %s\n", node.GetText(), node.GetReference()))
	// 		logWindow.ScrollToEnd()
	// 	})
	// 	treeView.SetDoneFunc(func(key tcell.Key) {
	// 		if key == tcell.KeyTab {
	// 			app.SetFocus(notes)
	// 		}
	// 	})
	// 	notes.SetDoneFunc(func() {
	// 		app.SetFocus(treeView)
	// 	})
	//
	// 	cmdBar.AddCmd(":mkbook", "make a top level notebook, prefix with // to place under current notebook", func(name string) error {
	// 		var err error
	// 		var parent *model.Notebook
	// 		if strings.HasPrefix(name, "//") {
	// 			parent, err = notebooks.NotebookByID(treeView.GetCurrentNode().GetReference().(string))
	// 			if err != nil {
	// 				logWindow.SetText(logWindow.GetText(false) + "couldn't find parent: " + err.Error())
	// 				logWindow.ScrollToEnd()
	// 				return err
	// 			}
	// 			logWindow.SetText(logWindow.GetText(false) + "adding notebook with parent " + parent.Title)
	// 			logWindow.ScrollToEnd()
	// 			name = name[2:]
	// 		}
	//
	// 		logWindow.SetText(logWindow.GetText(false) + "adding notebook " + name)
	// 		logWindow.ScrollToEnd()
	// 		commands <- Command{Execute: func(notebooks *model.NotebookService) error {
	// 			change.AddNotebook(name, parent).Apply(notebooks)
	// 			return nil
	// 		}}
	// 		cmdBar.SetText("")
	// 		app.SetFocus(treeView)
	//
	// 		return nil
	// 	})
	// 	cmdBar.AddCmd(":rmbook", "remove currently selected notebook", func(_ string) error {
	// 		id := treeView.GetCurrentNode().GetReference().(string)
	// 		logWindow.SetText(logWindow.GetText(false) + "deleting notebook " + id)
	// 		logWindow.ScrollToEnd()
	// 		commands <- Command{Execute: func(notebooks *model.NotebookService) error {
	// 			change.DeleteNotebook(id).Apply(notebooks)
	// 			return nil
	// 		}}
	// 		cmdBar.SetText("")
	// 		app.SetFocus(treeView)
	// 		return nil
	// 	})
	//
	// 	// build initial ui
	// 	for _, notebook := range notebooks.RootNotebooks() {
	// 		node := NotebookToTreeNode(notebook)
	//
	// 		treeView.GetRoot().AddChild(node)
	// 	}
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
