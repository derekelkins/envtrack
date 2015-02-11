package main

import (
	"encoding/json"
	"log"
	"net/url"

	"github.com/hashicorp/consul/api"
)

type ConsulKeyListener struct {
	client    *api.Client
	path      string
	waitIndex uint64
}

func NewConsulKeyListener(uri *url.URL) KeyListener {
	config := api.DefaultConfig()
	// TODO: Set WaitTime and various other things.
	if uri.Host != "" {
		config.Address = uri.Host
	}
	client, err := api.NewClient(config)
	if err != nil {
		log.Fatal("failed to connect to Consul: ", err)
	}
	return &ConsulKeyListener{client: client, path: uri.Path, waitIndex: 1}
}

// See http://godoc.org/github.com/hashicorp/consul/api#KV.List
func (l *ConsulKeyListener) GetKeys() ([]byte, error) {
	kv := l.client.KV()
	options := &api.QueryOptions{WaitIndex: l.waitIndex}
	kvps, qm, err := kv.List(l.path, options)
	if err != nil {
		return nil, err
	}
	l.waitIndex = qm.LastIndex
	return json.MarshalIndent(kvps, "", "\t")
}
