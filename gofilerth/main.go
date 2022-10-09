package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/harusame0616/GoFilerth/gofilerth/infrastructure/local"
	"github.com/harusame0616/GoFilerth/gofilerth/usecase"
)

func main() {

	fileQueryUsecase, _ := usecase.NewFileQuery(local.NewFileQuery())
	fmt.Print("一覧表示するパスを入力してください\n> ")

	var path string
	fmt.Scan(&path)

	files, err := fileQueryUsecase.ListFiles(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Printf(filepath.Join(file.Path, file.Name) + "\n")
	}
}
