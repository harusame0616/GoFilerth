package usecase

import (
	"fmt"
	"log"

	"github.com/harusame0616/GoFilerth/gofilerth/domain/file"
	"github.com/harusame0616/GoFilerth/gofilerth/domain/filer"
)

type CommandGateway interface {
	OpenShell(workDir string)
	OpenEditor(path string) error
}

type FilerCommand struct {
	fileRepository    file.Repository
	filerRepository   filer.Repository
	fileDomainService *filer.DomainService
	commandGateway    CommandGateway
}

func NewFilerCommand(fileRepository file.Repository, filerRepository filer.Repository, commandGateway CommandGateway) *FilerCommand {
	return &FilerCommand{
		fileRepository:    fileRepository,
		filerRepository:   filerRepository,
		fileDomainService: filer.NewDomainService(fileRepository),
		commandGateway:    commandGateway,
	}
}

func (filerCommand *FilerCommand) CreateNewFiler(initialPath string) (string, error) {
	if file, err := filerCommand.fileRepository.GetOneByPath(initialPath); err != nil {
		return "", err
	} else if !file.IsDirectory() {
		return "", fmt.Errorf("The initial path is not directory. (%s)", initialPath)
	} else {
		// do nothing
	}
	filer := filer.New(initialPath)

	filerCommand.filerRepository.Save(filer)
	return filer.Id(), nil
}

func (filerCommand *FilerCommand) ChangeDirectory(id string, destinationPath string) error {
	var filer *filer.Filer
	if _filer, err := filerCommand.filerRepository.GetOneById(id); err != nil {
		return err
	} else {
		filer = _filer
	}

	if err := filerCommand.fileDomainService.ChangeDirectory(filer, destinationPath); err != nil {
		return err
	}

	filerCommand.filerRepository.Save(filer)
	return nil
}

func (filerCommand *FilerCommand) UpDirectory(id string) (string, error) {
	filer, err := filerCommand.filerRepository.GetOneById(id)
	if err != nil {
		log.Fatal(err)
	}

	if err := filerCommand.fileDomainService.UpDirectory(filer); err != nil {
		return "", err
	}

	filerCommand.filerRepository.Save(filer)
	return filer.CurrentPath(), nil
}

func (filerCommand *FilerCommand) CurrentPath(id string) string {
	filer, err := filerCommand.filerRepository.GetOneById(id)
	if err != nil {
		log.Fatal(err)
	}

	return filer.CurrentPath()
}

func (filerCommand *FilerCommand) OpenShell(id string) {
	filer, err := filerCommand.filerRepository.GetOneById(id)
	if err != nil {
		log.Fatal(err)
	}

	filerCommand.commandGateway.OpenShell(filer.CurrentPath())
}

func (filerCommand *FilerCommand) OpenEditor(id string, path string) {
	filerCommand.commandGateway.OpenEditor(path)
}
