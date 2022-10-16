package shell

import (
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
