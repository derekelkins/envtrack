package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"
)

/* Examples for flag
var hostIp = flag.String("ip", "", "IP for ports mapped to the host")
var internal = flag.Bool("internal", false, "Use internal ports instead of published ones")
var refreshInterval = flag.Int("ttl-refresh", 0, "Frequency with which service TTLs are refreshed")
*/

type KeyListener interface {
    GetKeys() []byte, error
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
}

func NewKeyListener(uri *url.URL) KeyListener {
	factory := map[string]func(*url.URL) KeyListener{
		"consul":  NewConsulRegistry,
		//"etcd":    NewEtcdRegistry,
	}[uri.Scheme]
	if factory == nil {
		log.Fatal("unrecognized listener: ", uri.Scheme)
	}
	return factory(uri)
}
