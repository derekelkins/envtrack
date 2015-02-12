package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"os"

	"github.com/hashicorp/consul/api"
)

type ConsulKeyListener struct {
	client    *api.Client
	config    *api.Config
	path      string
	waitIndex uint64
}

func NewConsulKeyListener(uri *url.URL, disconnected bool) KeyListener {
	config := api.DefaultConfig()
	// TODO: Set WaitTime and various other things.
	if uri.Host != "" {
		config.Address = uri.Host
	}
	if disconnected {
		return &ConsulKeyListener{client: nil, config: config, path: uri.Path, waitIndex: 1}
	} else {
		client, err := api.NewClient(config)
		if err != nil {
			log.Fatal("failed to connect to Consul: ", err)
		}
		return &ConsulKeyListener{client: client, config: config, path: uri.Path, waitIndex: 1}
	}
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

func (l *ConsulKeyListener) Script(srcPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	kvps := make([]*api.KVPair, 0)

	dec := json.NewDecoder(src)
	err = dec.Decode(&kvps)
	if err != nil {
		return err
	}

	separator := rand.Int63()
	for _, kvp := range kvps {
		_, err := fmt.Printf("curl -X PUT '%s://%s/v1/kv/%s' <<EOF-%d\n%s\nEOF-%d\n",
			l.config.Scheme,
			l.config.Address,
			kvp.Key,
			separator,
			string(kvp.Value[:]),
			separator)
		if err != nil {
			return err
		}
	}

	return nil
}
