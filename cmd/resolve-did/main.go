package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/sfomuseum/go-atproto/plc"
)

func main() {

	var did string

	flag.StringVar(&did, "did", "", "The DID to resolve")
	flag.Parse()

	ctx := context.Background()

	d, err := plc.ResolveDID(ctx, did)

	if err != nil {
		log.Fatal(err)
	}

	err = d.Marshal(os.Stdout)

	if err != nil {
		log.Fatal(err)
	}
}
