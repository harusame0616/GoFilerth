package local

import (
	"os"

	"github.com/harusame0616/GoFilerth/gofilerth/usecase"
)

type FileQuery struct {
}

func NewFileQuery() *FileQuery {
	return &FileQuery{}
}

func (fileQuery *FileQuery) ListFiles(path string) ([]usecase.FileDto, error) {
	dirEntry, err := os.ReadDir(path)

	if err != nil {
		return nil, err
	}

	fileDtoList := make([]usecase.FileDto, len(dirEntry))
	for index, file := range dirEntry {
		fileInfo, err := file.Info()
		if err != nil {
			return nil, err
		}

		fileDtoList[index].Name = fileInfo.Name()
		fileDtoList[index].Path = path
		fileDtoList[index].Size = fileInfo.Size()
		fileDtoList[index].IsDirectory = fileInfo.IsDir()
		fileDtoList[index].ModifiedAt = fileInfo.ModTime()
	}

	return fileDtoList, nil
}
