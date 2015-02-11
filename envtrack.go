package main

import (
	"flag"
	"log"
	"net/url"
)

var backendFlag = flag.String("backend", "file", "Backend to use.  Either 'file' or 'git'.")
var backendPath = flag.String("path", "config", "Path to the file to store the saved key-value pairs.")
var scripting = flag.Bool("script", false, "Output script based on file pointed to by path.")

type KeyListener interface {
	GetKeys() ([]byte, error)
	Script(srcPath string) error
}

type Backend interface {
	Store(data []byte) error
}

func main() {
	flag.Parse()

	uri, err := url.Parse(flag.Arg(0))
	if err != nil {
		log.Fatal("envtrack: ", err)
	}

	listener := NewKeyListener(uri, *scripting)

	if *scripting {
		err := listener.Script(*backendPath)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		backend := NewBackend(*backendPath)

		// To allow requests to be received while we're writing to a file.
		pipe := make(chan []byte, 100)

		go func() {
			for {
				data, err := listener.GetKeys()
				if err != nil {
					log.Println(err)
					continue
				}
				pipe <- data
			}
		}()

		for data := range pipe {
			err := backend.Store(data)
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}
}

func NewKeyListener(uri *url.URL, scripting bool) KeyListener {
	factory := map[string]func(*url.URL, bool) KeyListener{
		"consul": NewConsulKeyListener,
		//"etcd":    NewEtcdKeyListener,
	}[uri.Scheme]
	if factory == nil {
		log.Fatal("unrecognized listener: ", uri.Scheme)
	}
	return factory(uri, scripting)
}

func NewBackend(path string) Backend {
	factory := map[string]func(string) Backend{
		"file": NewFileBackend,
		"git":  NewGitBackend,
	}[*backendFlag]
	if factory == nil {
		log.Fatal("unrecognized backend: ", *backendFlag)
	}
	return factory(path)
}
