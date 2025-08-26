package create

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var users_database_uri string

var handle string
var service string
var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("server")

	fs.StringVar(&users_database_uri, "user-database-uri", "mem://", "A valid gocloud.dev/blob.Bucket URI.")
	fs.StringVar(&handle, "handle", "", "The handle name for the new account.")
	fs.StringVar(&service, "service", "", "The service name for the new account.")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")
	return fs
}
