package filer

import (
	"strings"

	"github.com/google/uuid"
)

type Filer struct {
	id          string
	currentPath string
}

func New(initialPath string) *Filer {
	return &Filer{id: uuid.NewString(), currentPath: initialPath}
}

func (filer *Filer) Id() string {
	return filer.id
}

func (filer *Filer) CurrentPath() string {
	return filer.currentPath
}

func (_filer *Filer) changeDirecotry(destinationPath string) string {
	if strings.HasSuffix(destinationPath, "/") && len(destinationPath) > 1 {
		// ルート以外の時の末尾の/は除去する
		_filer.currentPath = destinationPath[:len(destinationPath)-1]
	} else {
		_filer.currentPath = destinationPath
	}

	return _filer.currentPath
}
