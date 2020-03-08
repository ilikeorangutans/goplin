package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/ilikeorangutans/goplin/pkg/cmdbar"
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
	app          *tview.Application
	cmdBar       *cmdbar.CmdBar
	treeView     *tview.TreeView
	log          *Log
	logger       Logger
	workspace    *workspace.Workspace
	commandQueue chan Command
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

	//treeView.SetCurrentNode(treeView.GetRoot().GetChildren()[0])

	notes := tview.NewList()
	notes.SetBorder(true).SetTitle("Notes")

	noteDetails := tview.NewFlex()

	mainPanel.AddItem(treeView, 0, 1, false)
	mainPanel.AddItem(notes, 0, 1, false)

	mainPanel.AddItem(noteDetails, 0, 2, false)

	root.AddItem(mainPanel, 0, 3, false)

	m.log = NewLog()
	m.logger = m.log
	m.cmdBar = cmdbar.New("begin entering commands by pressing \":\" or begin search with \"/\"")
	root.AddItem(m.log, 15, 1, false)
	root.AddItem(m.cmdBar, 1, 1, false)
	m.app.SetRoot(root, true)
	m.app.SetFocus(m.treeView)
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
	m.cmdBar.AddCmd(":foo", "foo", func(s string) error {
		m.logger.WriteString("doing foo")
		m.commandQueue <- &FooCommand{UIUpdater: m, Tree: m.treeView}
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
	if m.app.GetFocus() != m.cmdBar && event.Key() == tcell.KeyRune {
		switch event.Rune() {
		case ':':
			m.cmdBar.SetText(":")
			m.app.SetFocus(m.cmdBar)
		case '/':
			m.cmdBar.SetText("/")
			m.app.SetFocus(m.cmdBar)
		}
	}
	return event
}
