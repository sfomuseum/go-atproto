package pds

// https://atproto.com/specs/did
// https://github.com/did-method-plc/did-method-plc
// https://github.com/haileyok/cocoon/blob/main/plc/client.go#L57

import (
	"golang.org/x/crypto/ed25519"
)

func GenerateKeyPair() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	return ed25519.GenerateKey(nil)
}
