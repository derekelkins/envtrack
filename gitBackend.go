package main

import (
	"os"
	"os/exec"
	"time"
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

	err = exec.Command("git", "commit", "-m", time.Now().Format(time.RFC1123Z)).Run()
	return err
}
