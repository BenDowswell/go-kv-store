package main

import (
	"github.com/BenDowswell/go-kv-store/kvstore"
	"github.com/BenDowswell/go-kv-store/server"
)

func main() {
	kv := kvstore.NewKVStore()
	server.Start(kv)
}
