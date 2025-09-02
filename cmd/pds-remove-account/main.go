package main

import (
	"context"
	"log"

	_ "github.com/mattn/go-sqlite3"
	_ "gocloud.dev/blob/memblob"

	"github.com/sfomuseum/go-atproto/app/pds/account/delete"
	"github.com/sfomuseum/go-atproto/pds"
)

func main() {

	ctx := context.Background()

	err := pds.RegisterBlobAccountsSchemes(ctx)

	if err != nil {
		log.Fatalf("Failed to register blob schemes, %v", err)
	}

	err = delete.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to run remove account, %v", err)
	}
}
