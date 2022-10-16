package filer

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/harusame0616/GoFilerth/gofilerth/domain/file"
)

type DomainService struct {
	fileRepository file.Repository
}

func NewDomainService(fileRepository file.Repository) *DomainService {
	return &DomainService{fileRepository: fileRepository}
}

func (domainService *DomainService) ChangeDirectory(filer *Filer, destinationPath string) error {
	if strings.HasSuffix(destinationPath, "/") && destinationPath != "/" {
		destinationPath = destinationPath[:len(destinationPath)-1]
	}
	file, err := domainService.fileRepository.GetOneByPath(destinationPath)

	if err != nil {
		return err
	}

	if !file.IsDirectory() {
		return errors.New("The path is not directory.")
	}

	filer.changeDirecotry(destinationPath)

	return nil
}

func (domainService *DomainService) UpDirectory(filer *Filer) error {
	dir, file := filepath.Split(filer.currentPath)
	if file == "" {
		return errors.New("current path is root.")
	}
	return domainService.ChangeDirectory(filer, dir)
}
