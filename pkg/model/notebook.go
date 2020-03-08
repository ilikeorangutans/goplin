package model

type Notebook struct {
	Item
	Title     string
	Parent    *Notebook
	Order     int
	Notebooks []*Notebook
}

func (n Notebook) HasParent() bool {
	return n.Parent != nil
}

type Notebooks struct {
	Notebooks []*Notebook
	ByID      map[string]*Notebook
}
