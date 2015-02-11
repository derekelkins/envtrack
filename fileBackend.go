package main

import (
    "errors"
    "log"
    "os"
)

type FileBackend struct {
    path string
}

func NewFileBackend(path string) Backend {
    return &FileBackend{path: path}
}

func (b *FileBackend) Store(data []byte) error {

}
