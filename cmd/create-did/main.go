package main

import (
	"flag"
	"log"
	"os"

	"github.com/sfomuseum/go-atproto/did"
)

func main() {

	var host string
	var name string

	flag.StringVar(&name, "name", "alice", "The name of the account the DID is being created for.")
	flag.StringVar(&host, "host", "https://example.com", "The hostname for the account hosting {name}.")

	flag.Parse()

	d, _, err := did.NewDID(name, host)

	if err != nil {
		log.Fatalf("Failed to create DID, %v", err)
	}

	err = d.Marshal(os.Stdout)

	if err != nil {
		log.Fatalf("Failed to marshal DID, %v", err)
	}
}
