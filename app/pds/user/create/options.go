package create

import (
	"context"
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	UsersDatabaseURI string `json:"users_database_uri"`
}

func OptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	opts := &RunOptions{
		UsersDatabaseURI: users_database_uri,
	}

	return opts, nil
}
