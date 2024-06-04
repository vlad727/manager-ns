package main

import (
	"flag"
	"log"
	"manager-ns/api" // import api package with handlers.go to validate request
	"net/http"
)

const (
	port = ":8443"
)

var (
	tlscert, tlskey string
)

func main() {
	log.Printf("Application listening port %s\n", port)
	flag.StringVar(&tlscert, "tlsCertFile", "/certs/tls.crt",
		"File containing a certificate for HTTPS.")
	flag.StringVar(&tlskey, "tlsKeyFile", "/certs/tls.key",
		"File containing a private key for HTTPS.")
	flag.Parse()

	http.HandleFunc("/validate", api.Validate) // func validate in package api
	log.Fatal(http.ListenAndServeTLS(port, tlscert, tlskey, nil))
}
