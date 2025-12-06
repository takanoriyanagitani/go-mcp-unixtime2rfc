package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/takanoriyanagitani/go-mcp-unixtime2rfc"
)

const (
	defaultPort         = 12030
	readTimeoutSeconds  = 10
	writeTimeoutSeconds = 10
	maxHeaderExponent   = 20
)

var port = flag.Int("port", defaultPort, "port to listen")

func main() {
	flag.Parse()

	handler, err := unixtime2rfc.NewServer()
	if err != nil {
		log.Fatalf("failed to create server: %v\n", err)
	}

	address := fmt.Sprintf(":%d", *port)

	//nolint:exhaustruct
	server := &http.Server{
		Addr:           address,
		Handler:        handler,
		ReadTimeout:    readTimeoutSeconds * time.Second,
		WriteTimeout:   writeTimeoutSeconds * time.Second,
		MaxHeaderBytes: 1 << maxHeaderExponent,
	}

	log.Printf("ready to start http mcp server. listening on %s\n", address)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to listen and serve: %v\n", err)
	}
}
