package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/bluesky-social/indigo/atproto/crypto"
	"github.com/sfomuseum/go-atproto/plc"
)

func main() {

	var did string
	var cid string

	var mb_private string

	flag.StringVar(&did, "did", "", "...")
	flag.StringVar(&cid, "cid", "", "...")

	flag.StringVar(&mb_private, "private-key", "", "The private key used to sign the request encoded as a Multibase string.")
	flag.Parse()

	ctx := context.Background()

	// START OF put me in a function

	private_key, err := crypto.ParsePrivateMultibase(mb_private)

	if err != nil {
		log.Fatalf("Failed to derive private key, %v", err)
	}

	private_key_k256, err := crypto.ParsePrivateBytesK256(private_key.Bytes())

	if err != nil {
		log.Fatalf("Failed to create private key k256, %v", err)
	}

	// END OF put me in a function

	op, err := plc.TombstoneDID(ctx, did, cid, private_key_k256)

	if err != nil {
		log.Fatalf("Failed to create DID, %v", err)
	}

	fmt.Printf("DID tombstoned with CID '%s'\n", op.CID().String())
}
