package model

type Note struct {
	Item
	Title    string
	Body     string
	Notebook *Notebook
}

type Notebook struct {
	Item
	Title     string
	Parent    *Notebook
	Order     int
	Notebooks []*Notebook
}

type Notebooks struct {
	Notebooks []*Notebook
	ByID      map[string]*Notebook
}

// Roots returns all the root notebooks, i.e. notebooks without parents.
func (nb *Notebooks) Roots() []*Notebook {
	// TODO we should probably cache this

	var result []*Notebook
	for _, notebook := range nb.Notebooks {
		if notebook.ParentID == "" {
			result = append(result, notebook)
		}
	}

	return result
}
