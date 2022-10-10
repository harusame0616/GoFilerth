package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/harusame0616/GoFilerth/gofilerth/infrastructure/gcs"
	"github.com/harusame0616/GoFilerth/gofilerth/infrastructure/local"
	"github.com/harusame0616/GoFilerth/gofilerth/usecase"
)

func main() {
	fmt.Print("一覧表示するパスを入力してください\n> ")

	var path string
	fmt.Scan(&path)

	var fileQuery usecase.FileQueryInterface
	if strings.HasPrefix(path, "gcs://") {
		if gcsFileQuery, err := gcs.NewFileQueryByCredential("credentials/gcs.json"); err != nil {
			log.Fatal(err)
		} else {
			fileQuery = gcsFileQuery
		}
	} else {
		fileQuery = local.NewFileQuery()
	}

	fileQueryUsecase, err := usecase.NewFileQuery(fileQuery)
	files, err := fileQueryUsecase.ListFiles(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Printf(file.Path + "/" + filepath.Join(file.Name) + "\n")
	}
}
