package delete

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var database_uri string

var accounts_database_uri string
var keys_database_uri string
var operations_database_uri string

var did string
var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("server")

	fs.StringVar(&database_uri, "database-uri", "", "An optional common database URI to apply to all other empty -{SUBJECT}-database-uri flags. This is a convenience flag for things like SQL databases.")

	fs.StringVar(&accounts_database_uri, "account-database-uri", "", "A registered sfomuseum/go-atproto/pds.AccountsDatabase URI.")
	fs.StringVar(&keys_database_uri, "keys-database-uri", "", "A registered sfomuseum/go-atproto/pds.KeysDatabase URI.")
	fs.StringVar(&operations_database_uri, "operations-database-uri", "", "A registered sfomuseum/go-atproto/pds.OperationsDatabase URI.")

	fs.StringVar(&did, "did", "", "The DID for the account to delete.")

	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")
	return fs
}
