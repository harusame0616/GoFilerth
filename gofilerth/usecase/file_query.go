package usecase

import (
	"errors"
	"time"
)

type FileQueryInterface interface {
	ListFiles(path string) ([]FileDto, error)
}

type FileDto struct {
	Name       string
	Path       string
	Size       int64
	ModifiedAt time.Time
	FileType   string
}

type FileQueryUsecase struct {
	fileQueryInterface FileQueryInterface
}

func NewFileQuery(fileQuery FileQueryInterface) (*FileQueryUsecase, error) {
	if fileQuery == nil {
		return nil, errors.New("fileQuery is required.")
	}

	return &FileQueryUsecase{fileQueryInterface: fileQuery}, nil
}

func (fileQueryUsecase *FileQueryUsecase) ListFiles(path string) ([]FileDto, error) {
	return fileQueryUsecase.fileQueryInterface.ListFiles(path)
}
