package main

import (
	"os"
	"time"
)

type FileBackend struct {
	path string
}

func NewFileBackend(path string) Backend {
	return &FileBackend{path: path}
}

func (b *FileBackend) Store(data []byte) error {
	file, err := os.Create(b.path + "-" + time.Now().Format(time.RFC3339))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}
