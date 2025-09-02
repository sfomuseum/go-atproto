package list

import (
	"context"
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	AccountsDatabaseURI string `json:"accounts_database_uri"`
	Verbose             bool   `json:"verbose"`
}

func OptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	if database_uri != "" {

		if accounts_database_uri == "" {
			accounts_database_uri = database_uri
		}
	}

	opts := &RunOptions{
		AccountsDatabaseURI: accounts_database_uri,
		Verbose:             verbose,
	}

	return opts, nil
}
