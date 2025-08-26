package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/sfomuseum/go-atproto/plc"
)

func main() {

	var host string
	var handle string

	flag.StringVar(&handle, "handle", "alice", "The name of the account the DID is being created for.")
	flag.StringVar(&host, "host", "https://example.com", "The hostname for the account hosting {name}.")

	flag.Parse()

	ctx := context.Background()

	rsp, err := plc.NewDID(ctx, host, handle)

	if err != nil {
		log.Fatalf("Failed to create DID, %v", err)
	}

	did := rsp.DID

	err = did.Marshal(os.Stdout)

	if err != nil {
		log.Fatalf("Failed to marshal DID, %v", err)
	}
}
