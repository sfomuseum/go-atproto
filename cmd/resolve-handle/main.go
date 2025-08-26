package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/sfomuseum/go-atproto/plc"
)

func main() {

	var account_handle string
	var account_host string
	var newline bool

	flag.StringVar(&account_handle, "handle", "", "The ATProto handle to lookup.")
	flag.StringVar(&account_host, "host", "", "The ATProto hostname to query for the handle lookup.")
	flag.BoolVar(&newline, "with-newline", false, "Print final DID with trailing newline.")
	flag.Parse()

	ctx := context.Background()

	str_did, err := plc.ResolveHandle(ctx, account_handle, account_host)

	if err != nil {
		log.Fatal(err)
	}

	layout := "%s"

	if newline {
		layout = fmt.Sprintf("%s\n", layout)
	}

	fmt.Printf(layout, str_did)
}
