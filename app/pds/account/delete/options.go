package delete

import (
	"context"
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	AccountsDatabaseURI   string `json:"accounts_database_uri"`
	KeysDatabaseURI       string `json:"keys_database_uri"`
	OperationsDatabaseURI string `json:"operations_database_uri"`
	DID                   string `json:"did"`
}

func OptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	if database_uri != "" {

		if accounts_database_uri == "" {
			accounts_database_uri = database_uri
		}

		if keys_database_uri == "" {
			keys_database_uri = database_uri
		}

		if operations_database_uri == "" {
			operations_database_uri = database_uri
		}
	}

	opts := &RunOptions{
		AccountsDatabaseURI:   accounts_database_uri,
		KeysDatabaseURI:       keys_database_uri,
		OperationsDatabaseURI: operations_database_uri,
		DID:                   did,
	}

	return opts, nil
}
