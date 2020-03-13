package model

type Note struct {
	Item
	Title    string
	Body     string
	Notebook *Notebook
}
