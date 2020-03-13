package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/ilikeorangutans/goplin/pkg/cmdbar"
	"github.com/ilikeorangutans/goplin/pkg/model"
	"github.com/ilikeorangutans/goplin/pkg/workspace"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func NewMain() *Main {
	app := tview.NewApplication()
	main := &Main{
		app:          app,
		commandQueue: make(chan Command),
	}

	return main
}

type Main struct {
	app      *tview.Application
	cmdBar   *cmdbar.CmdBar
	treeView *tview.TreeView
	notes    *tview.List

	tabOrder   []tview.Primitive
	currentTab int

	log              *Log
	logger           Logger
	workspace        *workspace.Workspace
	commandQueue     chan Command
	SelectedNotebook *model.Notebook
}

// Run is the main entry point into the app.
func (m *Main) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	m.app.SetInputCapture(m.globalInputCapture)
	m.setUpUI()
	m.setUpCommands(ctx, cancel)

	m.logger.WriteString("goplin starting up...")

	if err := m.loadWorkspace(); err != nil {
		return err
	}

	go m.shutdownListener(ctx)
	go m.pollCommands(ctx)

	m.printHelp()

	return m.app.Run()
}

func (m *Main) pollCommands(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case cmd := <-m.commandQueue:
			err := cmd.Execute()
			// TODO can we split ui updates from the commands?
			if err != nil {
				panic(err)
			}
		}
	}
}

func (m *Main) loadWorkspace() error {
	// TODO here we'd actually read files
	m.logger.WriteString("loading workspace...")
	workspace, err := workspace.Open()
	if err != nil {
		return err
	}

	m.logger.WriteString(fmt.Sprintf("loaded workspace %q", workspace.Name()))

	m.workspace = workspace

	return nil
}

func (m *Main) setUpUI() {
	root := tview.NewFlex()
	root.SetDirection(tview.FlexRow)

	mainPanel := tview.NewFlex()

	treeView := tview.NewTreeView()
	treeView.SetTopLevel(1)
	treeView.SetBorder(true).SetTitle("Notebooks")
	m.treeView = treeView
	rootNode := tview.NewTreeNode("<root>")
	treeView.SetRoot(rootNode)
	treeView.SetSelectedFunc(func(node *tview.TreeNode) {
		notebook := node.GetReference().(*model.Notebook)
		m.SelectedNotebook = notebook
	})

	m.notes = tview.NewList()
	m.notes.SetBorder(true).SetTitle("Notes")

	noteDetails := tview.NewFlex()

	mainPanel.AddItem(treeView, 0, 1, false)
	mainPanel.AddItem(m.notes, 0, 1, false)

	mainPanel.AddItem(noteDetails, 0, 2, false)

	root.AddItem(mainPanel, 0, 3, false)

	m.log = NewLog()
	m.logger = m.log
	m.cmdBar = cmdbar.New("begin entering commands by pressing \":\" or begin search with \"/\"")
	root.AddItem(m.log, 15, 1, false)
	root.AddItem(m.cmdBar, 1, 1, false)
	m.app.SetRoot(root, true)
	m.app.SetFocus(m.treeView)

	m.tabOrder = []tview.Primitive{
		m.treeView,
		m.notes,
		//noteDetails,
	}
	m.treeView.SetDoneFunc(m.handleTabWithKey)
	m.notes.SetInputCapture(m.handleTabInputHandler)
}

func (m *Main) handleTabInputHandler(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyTab {
		m.handleTabWithKey(event.Key())
	}
	return event
}

func (m *Main) handleTabWithKey(key tcell.Key) {
	if key == tcell.KeyTab {
		m.currentTab = (m.currentTab + 1) % len(m.tabOrder)
		m.app.SetFocus(m.tabOrder[m.currentTab])
	}
}

func (m *Main) printHelp() {
	var sbuilder strings.Builder
	fmt.Fprint(&sbuilder, "Known commands:\n")
	for _, c := range m.cmdBar.SummarizeCommands() {
		fmt.Fprintf(&sbuilder, "%-10s  %s\n", c.Verb, c.Summary)
	}

	m.logger.WriteString(sbuilder.String())

}

func (m *Main) setUpCommands(ctx context.Context, cancel func()) {
	m.cmdBar.AddCmd(":q", "quit goplin", func(_ string) error {
		cancel()
		return nil
	})
	m.cmdBar.AddCmd(":help", "shows all known commands and their usage", func(_ string) error {
		m.printHelp()
		return nil
	})
	m.cmdBar.AddCmd(":mkbook", "create new notebook", func(s string) error {
		m.commandQueue <- &CreateNotebookCommand{UIUpdater: m, Tree: m.treeView, Workspace: m.workspace, Name: s}
		return nil
	})
	m.cmdBar.AddCmd(":mknote", "create new note in the current notebook", func(s string) error {
		currentNode := m.treeView.GetCurrentNode()
		if currentNode == nil {
			return fmt.Errorf("cannot create note without selecting a notebook first")
		}
		notebook := currentNode.GetReference().(*model.Notebook)
		m.commandQueue <- &CreateNoteCommand{UIUpdater: m, Notes: m.notes, Workspace: m.workspace, Name: s, Parent: notebook}
		return nil
	})
}

func (m *Main) QueueUpdateDraw(f func()) {
	m.app.QueueUpdateDraw(f)
}

// shutdownListener listens to the current context and stops the app when the context is done.
func (m *Main) shutdownListener(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			m.app.Stop()
		}
	}
}

func (m *Main) globalInputCapture(event *tcell.EventKey) *tcell.EventKey {
	if m.app.GetFocus() == m.cmdBar {
		if event.Key() == tcell.KeyEscape {
			m.app.SetFocus(m.tabOrder[m.currentTab])
		}
	} else {
		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case ':':
				m.cmdBar.SetText(":")
				m.app.SetFocus(m.cmdBar)
			case '/':
				m.cmdBar.SetText("/")
				m.app.SetFocus(m.cmdBar)
			}
		}
	}

	return event
}
