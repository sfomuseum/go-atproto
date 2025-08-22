package server

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var verbose bool
var server_uri string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("server")

	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")
	fs.StringVar(&server_uri, "server-uri", "http://localhost:8080", "A valid aaronland/go-http/v3/server.Server URI.")

	return fs
}
