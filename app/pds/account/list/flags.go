package list

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var database_uri string

var accounts_database_uri string
var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("server")

	fs.StringVar(&database_uri, "database-uri", "", "An optional common database URI to apply to all other empty -{SUBJECT}-database-uri flags. This is a convenience flag for things like SQL databases.")

	fs.StringVar(&accounts_database_uri, "account-database-uri", "", "A registered sfomuseum/go-atproto/pds.AccountsDatabase URI.")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")
	return fs
}
