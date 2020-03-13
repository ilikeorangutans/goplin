package workspace

func Open() (*Workspace, error) {
	return &Workspace{
		notebooks: NewNotebooks(),
		notes:     NewNotes(),
	}, nil
}

type Workspace struct {
	notebooks *Notebooks
	notes     *Notes
}

func (w *Workspace) Name() string {
	return "default workspace"
}

func (w *Workspace) Notebooks() *Notebooks {
	return w.notebooks
}

func (w *Workspace) Notes() *Notes {
	return w.notes
}
