package cmdbar

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func New(placeholder string) *CmdBar {
	inputField := tview.NewInputField()
	inputField.SetPlaceholder(placeholder)

	cmdbar := &CmdBar{
		InputField: inputField,
		commands:   make(map[string]*Cmd),
	}

	cmdbar.setupEventHandlers()

	return cmdbar
}

type CmdBar struct {
	*tview.InputField
	commands map[string]*Cmd
}

func (c *CmdBar) setupEventHandlers() {
	c.SetDoneFunc(c.doneFunc)
}

func (c *CmdBar) doneFunc(key tcell.Key) {
	if key != tcell.KeyEnter {
		c.Clear()
		return
	}

	input := strings.TrimSpace(c.GetText())
	if input == "" {
		return
	}

	parts := strings.SplitN(input, " ", 2)
	verb := parts[0]
	rest := ""
	if len(parts) > 1 {
		rest = parts[1]
	}

	cmd := c.FindByVerb(verb)
	if cmd == nil {
		c.Clear()
		c.SetPlaceholder(fmt.Sprintf("unknown verb %s", verb))
	} else {
		// TODO handle returned error
		cmd.Execute(rest)
		c.Clear()
	}
}

func (c *CmdBar) FindByVerb(verb string) *Cmd {
	if cmd, ok := c.commands[verb]; ok {
		return cmd
	} else {
		return nil
	}
}

// Clear clears the CmdBar
func (c *CmdBar) Clear() {
	c.SetText("")
}

func (c *CmdBar) AddCmd(verb, summary string, f func(string) error) {
	c.addCmd(&Cmd{
		CmdSummary: CmdSummary{
			Verb:    verb,
			Summary: summary,
		},
		Execute: f,
	})
}

func (c *CmdBar) addCmd(cmd *Cmd) {
	c.commands[cmd.Verb] = cmd
}

func (c *CmdBar) SummarizeCommands() []CmdSummary {
	summary := make([]CmdSummary, 0, len(c.commands))
	for _, cmd := range c.commands {
		summary = append(summary, cmd.CmdSummary)
	}
	return summary
}

type CmdSummary struct {
	Verb    string
	Summary string
}

type Cmd struct {
	CmdSummary
	Execute func(string) error
}
