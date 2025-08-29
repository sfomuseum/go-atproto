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

	opts := &RunOptions{
		AccountsDatabaseURI:   accounts_database_uri,
		KeyPairsDatabaseURI:   keypairs_database_uri,
		OperationsDatabaseURI: operations_database_uri,
		Handle:                handle,
		Service:               service,
	}

	return opts, nil
}
