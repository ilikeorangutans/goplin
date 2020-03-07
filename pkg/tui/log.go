package tui

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

type Logger interface {
	WriteString(string)
}

func NewLog() *Log {
	textview := tview.NewTextView()
	textview.SetScrollable(true).SetBorder(true).SetTitle("log")
	return &Log{
		TextView: textview,
	}
}

type Log struct {
	*tview.TextView
	builder strings.Builder
}

func (l *Log) WriteString(s string) {
	if s == "" {
		return
	}
	// TODO here we should probably throw away old lines
	fmt.Fprintln(l, s)
	l.ScrollToEnd()
}
