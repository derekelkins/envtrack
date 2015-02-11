package main

import (
	"flag"
	"log"
	"net/url"
)

/* Examples for flag
var hostIp = flag.String("ip", "", "IP for ports mapped to the host")
var internal = flag.Bool("internal", false, "Use internal ports instead of published ones")
var refreshInterval = flag.Int("ttl-refresh", 0, "Frequency with which service TTLs are refreshed")
*/

type KeyListener interface {
	GetKeys() ([]byte, error)
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
	listener := NewKeyListener(uri)

	path := "TODO"
	backend := NewBackend(path)

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

func NewKeyListener(uri *url.URL) KeyListener {
	factory := map[string]func(*url.URL) KeyListener{
		"consul": NewConsulKeyListener,
		//"etcd":    NewEtcdRegistry,
	}[uri.Scheme]
	if factory == nil {
		log.Fatal("unrecognized listener: ", uri.Scheme)
	}
	return factory(uri)
}

func NewBackend(path string) Backend {
	return NewFileBackend(path)
}
