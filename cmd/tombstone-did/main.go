package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/sfomuseum/go-atproto/crypto"
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

	private_key_k256, err := crypto.PrivateKeyK256FromMultibase(mb_private)

	if err != nil {
		log.Fatalf("Failed to create private key k256, %v", err)
	}

	plc_cl := plc.DefaultClient()

	op, err := plc.TombstoneDID(ctx, plc_cl, did, cid, private_key_k256)

	if err != nil {
		log.Fatalf("Failed to create DID, %v", err)
	}

	fmt.Printf("DID tombstoned with CID '%s'\n", op.CID().String())
}
