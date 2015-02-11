package main

import (
	"os"
	"os/exec"
)

type GitBackend struct {
	path string
}

func NewGitBackend(path string) Backend {
	return &GitBackend{path: path}
}

func (b *GitBackend) Store(data []byte) error {
	file, err := os.Create(b.path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	err = exec.Command("git", "add", b.path).Run()
	if err != nil {
		return err
	}
	return exec.Command("git", "commit", "-m", "envtrack auto-commit").Run()
}
