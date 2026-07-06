package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/BenDowswell/go-kv-store/kvstore"
    "github.com/BenDowswell/go-kv-store/server"
)

func main() {

	httpAddr := flag.String("http", ":8080", "HTTP listen address")
	tcpAddr := flag.String("tcp", ":9090", "TCP listen address")
	flag.Parse()

	// start up store
	kv := kvstore.NewKVStore()

	// start http passing in httpaddr flag
	go func() {
		log.Printf("HTTP listening on %s", *httpAddr)
		if err := server.Start(kv, *httpAddr); err != nil {
			log.Fatalf("HTTP: %v", err)
		}
	}()
	// start tcp passing in tcpaddr flag
	go func() {
		log.Printf("TCP listening on %s", *tcpAddr)
		if err := server.ServeTCP(kv, *tcpAddr); err != nil {
			log.Fatalf("TCP: %v", err)
		}
	}()
	// listen out for exit and signal shutting down gracefully
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("shutting down")
}
