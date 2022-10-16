package shell

import (
	"errors"
	"os"
	"os/exec"
)

type CommandGateway struct {
}

func NewCommandGateway() *CommandGateway {
	return &CommandGateway{}
}

func (_ *CommandGateway) OpenShell(workDir string) {
	defaultShell := os.Getenv("SHELL")
	shell := exec.Command(defaultShell)
	shell.Stdin = os.Stdin
	shell.Stdout = os.Stdout
	shell.Stderr = os.Stderr
	shell.Dir = workDir
	shell.Run()
}

func (_ *CommandGateway) OpenEditor(path string) error {
	defaultEditor := os.Getenv("EDITOR")
	if defaultEditor == "" {
		return errors.New("Please set $EDITOR Environment variable")
	}

	shell := exec.Command(defaultEditor, path)
	shell.Stdin = os.Stdin
	shell.Stdout = os.Stdout
	shell.Stderr = os.Stderr
	shell.Run()
	return nil
}
