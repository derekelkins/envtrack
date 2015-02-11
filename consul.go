package main

import (
	"errors"
	"fmt"
	"log"
    "net/url"

    "github.com/hashicorp/consul/api"
)

type ConsulKeyListener struct {
    client api.Client
    path string
}

func NewConsulKeyListener(uri *url.URL) KeyListener {
    config := api.DefaultConfig()
	if uri.Host != "" {
		config.Address = uri.Host
	}
    client, err := api.NewClient(config)
    if err != nil {
        log.Fatal("failed to connect to Consul: ", err)
    }
	return &ConsulKeyListener{client: client, path: uri.Path}
}

func (l *ConsulKeyListener) GetKeys() []byte, error {
    kv := client.KV()
// Use KV.List http://godoc.org/github.com/hashicorp/consul/api#KV.List

}
