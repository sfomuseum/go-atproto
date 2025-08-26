package plc

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/base64"
	"fmt"

	"github.com/bluesky-social/indigo/atproto/identity"
	"github.com/bluesky-social/indigo/atproto/syntax"
	"github.com/fxamacker/cbor/v2"
	"github.com/whyrusleeping/go-did"
)

// https://atproto.com/specs/did
// https://atproto.com/specs/handle#handle-resolution

// For reference, this is what Blacksky does:
// https://github.com/blacksky-algorithms/rsky/blob/main/rsky-pds/src/plc/operations.rs#L50C5-L58C46
// https://github.com/blacksky-algorithms/rsky/blob/main/rsky-pds/src/plc/operations.rs#L281
// https://github.com/blacksky-algorithms/rsky/blob/main/rsky-common/src/sign.rs#L8

type NewDIDResult struct {
	DID          *identity.DIDDocument
	PlcOperation PlcOperationSigned
	PrivateKey   *did.PrivKey
}

// https://github.com/bluesky-social/indigo/blob/main/atproto/identity/identity.go#L42	<-- ParseIdentity (from DIDDoc)
// https://github.com/bluesky-social/indigo/blob/8be102876fb7e638aa4c9ed6c9d4991ca19a0973/atproto/identity/diddoc.go#L7	<-- DIDDocument

func NewDID(ctx context.Context, service string, handle string) (*NewDIDResult, error) {

	// https://web.plc.directory/spec/v0.1/did-plc
	// In pseudo-code: did:plc:${base32Encode(sha256(createOp)).slice(0,24)}

	// Collect values for the essential operation data fields, including generating new secure key pairs if necessary
	// Only secp256k1 (“k256”) and NIST P-256 (“p256”) keys are currently supported for rotation keys, whereas verificationMethods keys can be any syntactically-valid did:key.

	private_key, err := did.GeneratePrivKey(rand.Reader, did.KeyTypeP256)

	if err != nil {
		return nil, fmt.Errorf("Failed to generate private key, %w", err)
	}

	public_key := private_key.Public()
	public_mb := public_key.MultibaseString()

	verification_key := fmt.Sprintf("%s:%s", DID_KEY, public_mb)

	// Construct an “unsigned” regular operation object.
	// Include a prev field with null value. do not use the deprecated/legacy operation format for new DID creations

	unsigned_op := PlcOperation{
		Type: "plc_operation",
		VerificationMethods: map[string]string{
			"atproto": verification_key,
		},
		RotationKeys: []string{
			verification_key,
		},
		AlsoKnownAs: []string{fmt.Sprintf("at://%s", handle)},
		Services: map[string]PlcOperationService{
			"atproto_pds": {
				Type:     "AtprotoPersonalDataServer",
				Endpoint: service,
			},
		},
		Prev: nil, // genesis – no previous (Prev) operation
	}

	// Serialize the “unsigned” operation with DAG-CBOR, and sign the resulting bytes with one of the initial rotationKeys.
	// Encode the signature as base64url, and use that to construct a “signed” operation object

	enc_opts := cbor.CanonicalEncOptions()

	enc_mode, err := enc_opts.EncMode()

	if err != nil {
		return nil, fmt.Errorf("cbor encoder: %w", err)
	}

	unsigned_b, err := enc_mode.Marshal(unsigned_op)

	if err != nil {
		return nil, fmt.Errorf("unsigned CBOR marshal: %w", err)
	}

	sig, err := private_key.Sign(unsigned_b)

	sig_b64 := base64.RawURLEncoding.EncodeToString(sig)

	// Serialize the “signed” operation with DAG-CBOR, take the SHA-256 hash of those bytes, and encode the hash bytes in base32.
	// use the first 24 characters to generate DID value (did:plc:<hashchars>)

	signed_op := PlcOperationSigned{
		PlcOperation: unsigned_op,
		Signature:    sig_b64,
	}

	signed_b, err := enc_mode.Marshal(signed_op)

	if err != nil {
		return nil, fmt.Errorf("signed CBOR marshal: %w", err)
	}

	hash := sha256.Sum256(signed_b)

	b32_enc := base32.StdEncoding.WithPadding(base32.NoPadding)
	hash_b32 := b32_enc.EncodeToString(hash[:]) // 52 chars

	if len(hash_b32) < 24 {
		return nil, fmt.Errorf("hash too short")
	}

	str_did := hash_b32[:24]

	did_id := fmt.Sprintf("%s:%s", DID_PLC, str_did)

	_, err = did.ParseDID(did_id)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse plc did, %w", err)
	}

	// Update to use indigo/atproto/identity.DidDocument

	did := &identity.DIDDocument{
		DID:         syntax.DID(did_id),
		AlsoKnownAs: unsigned_op.AlsoKnownAs,
		VerificationMethod: []identity.DocVerificationMethod{
			identity.DocVerificationMethod{
				ID:                 fmt.Sprintf("%s#atproto", did_id),
				Type:               "Multikey",
				Controller:         did_id,
				PublicKeyMultibase: verification_key,
			},
		},
		Service: []identity.DocService{
			identity.DocService{
				ID:              "#atproto_pds",
				Type:            "AtprotoPersonalDataServer",
				ServiceEndpoint: service,
			},
		},
	}

	//

	rsp := &NewDIDResult{
		DID:          did,
		PlcOperation: signed_op,
		PrivateKey:   private_key,
	}

	return rsp, nil
}
