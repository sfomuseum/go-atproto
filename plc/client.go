package plc

import (
	"github.com/did-method-plc/go-didplc"
)

func DefaultClient() *didplc.Client {

	return &didplc.Client{
		DirectoryURL: "https://plc.directory",
		UserAgent:    "sfomuseum/go-atproto",
	}

}
