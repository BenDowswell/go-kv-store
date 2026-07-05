package main

import (
	"log"

	"github.com/BenDowswell/go-kv-store/kvstore"
	"github.com/BenDowswell/go-kv-store/server"
)

func main() {
	kv := kvstore.NewKVStore()

	if err := server.Start(kv); err != nil {
		log.Fatal(err)
	}
}
