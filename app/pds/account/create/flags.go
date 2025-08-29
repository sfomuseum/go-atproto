package create

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var accounts_database_uri string
var keypairs_database_uri string
var operations_database_uri string

var handle string
var service string
var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("server")

	fs.StringVar(&accounts_database_uri, "account-database-uri", "", "A registered sfomuseum/go-atproto/pds.AccountsDatabase URI.")
	fs.StringVar(&keypairs_database_uri, "keypairs-database-uri", "", "A registered sfomuseum/go-atproto/pds.KeyPairsDatabase URI.")
	fs.StringVar(&operations_database_uri, "operations-database-uri", "", "A registered sfomuseum/go-atproto/pds.OperationsDatabase URI.")

	fs.StringVar(&handle, "handle", "", "The handle name for the new account.")
	fs.StringVar(&service, "service", "", "The service name for the new account.")

	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")
	return fs
}
