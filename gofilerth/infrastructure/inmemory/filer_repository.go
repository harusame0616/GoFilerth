package inmemory

import (
	"errors"

	"github.com/harusame0616/GoFilerth/gofilerth/domain/filer"
)

type FilerRepository struct {
}

var filerStore = make(map[string]*filer.Filer, 2)

func NewFilerRepository() *FilerRepository {
	return &FilerRepository{}
}

func (repository *FilerRepository) GetOneById(id string) (*filer.Filer, error) {
	filer, ok := filerStore[id]
	if !ok {
		return nil, errors.New("the filer is not found.")
	}

	return filer, nil
}

func (repository *FilerRepository) Save(filer *filer.Filer) error {
	filerStore[filer.Id()] = filer

	return nil
}
