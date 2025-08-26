package main

import (
	"context"
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

	d, err := api.ResolveDID(ctx, did)

	if err != nil {
		log.Fatal(err)
	}

	err = d.Marshal(os.Stdout)

	if err != nil {
		log.Fatal(err)
	}
}
