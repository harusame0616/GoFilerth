package local

import (
	"os"
	"path/filepath"
	"strings"

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
		fileDtoList[index].FullPath = filepath.Join(path, fileInfo.Name())
		fileDtoList[index].Size = fileInfo.Size()
		fileDtoList[index].ModifiedAt = fileInfo.ModTime()

		switch {
		case fileInfo.IsDir():
			fileDtoList[index].FileType = "_directory"
		default:
			// ファイルタイプ判定
			splits := strings.Split(fileInfo.Name(), ".")
			if len(splits) <= 1 {
				fileDtoList[index].FileType = ""
			} else {
				fileDtoList[index].FileType = splits[len(splits)-1]
			}
		}
	}

	return fileDtoList, nil
}
