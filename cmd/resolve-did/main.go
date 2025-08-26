package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/sfomuseum/go-atproto/plc/api"
)

func main() {

	var did string

	flag.StringVar(&did, "did", "", "The DID to resolve")
	flag.Parse()

	ctx := context.Background()

	doc, err := api.ResolveDID(ctx, did)

	if err != nil {
		log.Fatal(err)
	}

	enc := json.NewEncoder(os.Stdout)
	err = enc.Encode(doc)

	if err != nil {
		log.Fatal(err)
	}
}
