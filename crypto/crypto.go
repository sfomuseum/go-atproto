package crypto

import (
	at_crypto "github.com/bluesky-social/indigo/atproto/crypto"
)

func PrivateKeyK256FromMultibase(mb_key string) (*at_crypto.PrivateKeyK256, error) {

	private_key, err := at_crypto.ParsePrivateMultibase(mb_key)

	if err != nil {
		return nil, err
	}

	return at_crypto.ParsePrivateBytesK256(private_key.Bytes())
}
