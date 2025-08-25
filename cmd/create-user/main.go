package main

import (
	"context"
	"log"

	_ "gocloud.dev/blob/memblob"

	"github.com/sfomuseum/go-atproto/app/pds/user/create"
	"github.com/sfomuseum/go-atproto/pds"
)

func main() {

	ctx := context.Background()

	err := pds.RegisterBlobUsersSchemes(ctx)

	if err != nil {
		log.Fatal(err)
	}

	err = create.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to run create user, %v", err)
	}
}
