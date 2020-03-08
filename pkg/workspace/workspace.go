package workspace

import "github.com/ilikeorangutans/goplin/pkg/model"

func Open() (*Workspace, error) {
	return &Workspace{
		notebookService: model.NewNotebookService(),
	}, nil
}

type Workspace struct {
	notebookService *model.NotebookService
}

func (w *Workspace) Name() string {
	return "default workspace"
}

func (w *Workspace) Notebooks() *model.NotebookService {
	return w.notebookService
}
