package create

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var verbose bool
var users_database_uri string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("server")

	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")
	fs.StringVar(&users_database_uri, "user-database-uri", "mem://", "A valid gocloud.dev/blob.Bucket URI.")

	return fs
}
