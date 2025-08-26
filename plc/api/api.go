package api

// https://web.plc.directory/
// https://web.plc.directory/api/redoc
// https://web.plc.directory/api/plc-server-openapi3.yaml

// Seemingly only GET methods...
// https://github.com/did-method-plc/did-method-plc/blob/main/website/server.go
// POST stuff (and other GET things) are over here...
// https://github.com/did-method-plc/did-method-plc/blob/main/packages/server

import (
	"net/url"
)

// Constant for the "plc.directory" host.
const PLC_DIRECTORY string = "plc.directory"

// Create a new `url.URL` instance for the "plc.directory" host no path, header or query details.
func NewURL() *url.URL {

	u := new(url.URL)
	u.Scheme = "https"
	u.Host = PLC_DIRECTORY

	return u
}
