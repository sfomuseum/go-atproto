package pds

// https://github.com/did-method-plc/did-method-plc

import (
	"golang.org/x/crypto/ed25519"
)

func GenerateKeyPair() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	return ed25519.GenerateKey(nil)
}
