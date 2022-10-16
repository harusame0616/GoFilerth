package local

import (
	"os"

	"github.com/harusame0616/GoFilerth/gofilerth/domain/file"
)

type FileRepository struct {
}

func NewFileRepository() *FileRepository {
	return &FileRepository{}
}

func (repository *FileRepository) GetOneByPath(path string) (*file.File, error) {
	if fileInfo, err := os.Stat(path); err == nil {
		return file.FromDto(file.Dto{
			Path:        path,
			IsDirectory: fileInfo.IsDir(),
		}), nil
	} else {
		return nil, err
	}
}
