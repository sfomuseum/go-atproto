package api

// https://web.plc.directory/
// https://web.plc.directory/api/redoc

import (
	"net/url"
)

const PLC_DIRECTORY string = "plc.directory"

func NewURL() *url.URL {

	u := new(url.URL)
	u.Scheme = "https"
	u.Host = PLC_DIRECTORY

	return u
}
