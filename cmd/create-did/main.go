package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/sfomuseum/go-atproto/plc"
)

func main() {

	var service string
	var handle string

	flag.StringVar(&handle, "handle", "alice", "The name of the account the DID is being created for.")
	flag.StringVar(&service, "service", "https://example.com", "The servicename for the account serviceing {name}.")

	flag.Parse()

	ctx := context.Background()

	plc_cl := plc.DefaultClient()

	rsp, err := plc.NewDID(ctx, plc_cl, service, handle)

	if err != nil {
		log.Fatalf("Failed to create DID, %v", err)
	}

	did := rsp.DID

	enc := json.NewEncoder(os.Stdout)
	err = enc.Encode(did)

	if err != nil {
		log.Fatalf("Failed to marshal DID, %v", err)
	}
}
