package create

import (
	"context"
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	AccountsDatabaseURI   string `json:"accounts_database_uri"`
	KeyPairsDatabaseURI   string `json:"keypairs_database_uri"`
	OperationsDatabaseURI string `json:"operations_database_uri"`
	Handle                string `json:"handle"`
	Service               string `json:"service"`
}

func OptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	if database_uri != "" {

		if accounts_database_uri == "" {
			accounts_database_uri = database_uri
		}

		if keypairs_database_uri == "" {
			keypairs_database_uri = database_uri
		}

		if operations_database_uri == "" {
			operations_database_uri = database_uri
		}
	}

	opts := &RunOptions{
		AccountsDatabaseURI:   accounts_database_uri,
		KeyPairsDatabaseURI:   keypairs_database_uri,
		OperationsDatabaseURI: operations_database_uri,
		Handle:                handle,
		Service:               service,
	}

	return opts, nil
}
