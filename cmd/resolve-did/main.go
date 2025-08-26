package main

// go run cmd/resolve-handle/main.go -service https://bsky.social -handle {ACCOUNT}.bsky.social | go run cmd/resolve-did/main.go -stdin | jq

import (
	"context"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
	"strings"

	"github.com/sfomuseum/go-atproto/plc/api"
)

func main() {

	var did string
	var stdin bool

	flag.StringVar(&did, "did", "", "The DID to resolve.")
	flag.BoolVar(&stdin, "stdin", false, "If true read DID from STDIN.")

	flag.Parse()

	ctx := context.Background()

	if stdin {

		b, err := io.ReadAll(os.Stdin)

		if err != nil {
			log.Fatal(err)
		}

		did = strings.TrimSpace(string(b))
	}

	if did == "" {
		log.Fatalf("Missing DID")
	}

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
