package main

import (
	"context"
	"log"

	"github.com/sfomuseum/go-pds/app/user/create"
)

func main() {

	ctx := context.Background()
	err := create.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to run create user, %v", err)
	}
}
