package plc

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	"github.com/fxamacker/cbor/v2"
	"github.com/multiformats/go-multibase"
)

// https://atproto.com/specs/did
// https://atproto.com/specs/handle#handle-resolution

// For reference, this is what Blacksky does:
// https://github.com/blacksky-algorithms/rsky/blob/main/rsky-pds/src/plc/operations.rs#L50C5-L58C46
// https://github.com/blacksky-algorithms/rsky/blob/main/rsky-pds/src/plc/operations.rs#L281
// https://github.com/blacksky-algorithms/rsky/blob/main/rsky-common/src/sign.rs#L8

type VerificationMethod struct {
	// The DID followed by an identifying fragment. Use #atproto as the fragment for atproto signing keys
	Id string `json:"id"`
	// The fixed string Multikey
	Type string `json:"type"`
	// DID controlling the key, which in the current version of atproto must match the account DID itself
	Controller string `json:"controller"`
	// The public key itself, encoded in multibase format (with multicodec type indicator, and "compressed" key bytes)
	PublicKeyMultibase string `json:"publicKeyMultibase"`
}

type Service struct {
	Id              string `json:"id"`
	Type            string `json:"type"`
	ServiceEndpoint string `json:"serviceEndpoint"`
}

type DID struct {
	Context             any                   `json:"@context"`
	Id                  string                `json:"id"`
	VerificationMethods []*VerificationMethod `json:"verificationMethods"`
	AlsoKnownAs         []string              `json:"alsoKnownAs"`
	Service             []*Service            `json:"service"`
}

func (d *DID) String() string {
	return d.Id
}

func (d *DID) Marshal(wr io.Writer) error {
	enc := json.NewEncoder(wr)
	return enc.Encode(d)
}

type NewDIDResult struct {
	DID             *DID
	CreateOperation CreatePlcOperationSigned
	PrivateKey      *ecdsa.PrivateKey
}

// https://github.com/bluesky-social/indigo/blob/main/atproto/identity/identity.go#L42	<-- ParseIdentity (from DIDDoc)
// https://github.com/bluesky-social/indigo/blob/8be102876fb7e638aa4c9ed6c9d4991ca19a0973/atproto/identity/diddoc.go#L7	<-- DIDDocument

func NewDID(ctx context.Context, host string, handle string) (*NewDIDResult, error) {

	// https://web.plc.directory/spec/v0.1/did-plc
	// In pseudo-code: did:plc:${base32Encode(sha256(createOp)).slice(0,24)}

	// Collect values for the essential operation data fields, including generating new secure key pairs if necessary
	// Only secp256k1 (“k256”) and NIST P-256 (“p256”) keys are currently supported for rotation keys, whereas verificationMethods keys can be any syntactically-valid did:key.

	private_key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	if err != nil {
		return nil, fmt.Errorf("Failed to create P256 key, %w", err)
	}

	public_key := &private_key.PublicKey

	public_key_b, err := public_key.Bytes()

	if err != nil {
		return nil, fmt.Errorf("Failed to derive publi key bytes", err)
	}

	prefix_mb := []byte{0x80, 0x12}
	data_mb := append(prefix_mb, public_key_b...)

	public_mb, err := multibase.Encode(multibase.Base58BTC, data_mb)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive multibase encoding for public key, %w", err)
	}

	verification_key := fmt.Sprintf("%s:%s", DID_KEY, public_mb)

	// Construct an “unsigned” regular operation object.
	// Include a prev field with null value. do not use the deprecated/legacy operation format for new DID creations

	unsigned_op := CreatePlcOperation{
		Type: "plc_operation",
		VerificationMethods: map[string]string{
			"atproto": verification_key,
		},
		RotationKeys: []string{
			verification_key,
		},
		AlsoKnownAs: []string{fmt.Sprintf("at://%s", handle)},
		Services: map[string]CreatePlcService{
			"atproto_pds": {
				Type:     "AtprotoPersonalDataServer",
				Endpoint: host,
			},
		},
		// genesis – no previous (Prev) operation
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

	// hash unsigned_b because that is what blacksky does...
	// https://github.com/blacksky-algorithms/rsky/blob/main/rsky-common/src/sign.rs#L8

	unsigned_hash := sha256.Sum256(unsigned_b)
	sig, err := ecdsa.SignASN1(rand.Reader, private_key, unsigned_hash[:])

	if err != nil {
		return nil, err
	}

	sig_b64 := base64.RawURLEncoding.EncodeToString(sig)

	// Serialize the “signed” operation with DAG-CBOR, take the SHA-256 hash of those bytes, and encode the hash bytes in base32.
	// use the first 24 characters to generate DID value (did:plc:<hashchars>)

	signed_op := CreatePlcOperationSigned{
		CreatePlcOperation: unsigned_op,
		Signature:          sig_b64,
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

	did := &DID{
		Context: []string{
			"https://www.w3.org/ns/did/v1",
			"https://w3id.org/security/multikey/v1",
			"https://w3id.org/security/suites/secp256k1-2019/v1",
		},
		Id: did_id,
		VerificationMethods: []*VerificationMethod{
			&VerificationMethod{
				Id:                 fmt.Sprintf("%s#atproto", did_id),
				Type:               "Multikey",
				Controller:         did_id,
				PublicKeyMultibase: fmt.Sprintf("%s:%s", DID_KEY, public_mb),
			},
		},
		// RotationKeys: unsigned.RotationKeys,
		AlsoKnownAs: unsigned_op.AlsoKnownAs,
		Service: []*Service{
			&Service{
				Id:              "#atproto_pds",
				Type:            "AtprotoPersonalDataServer",
				ServiceEndpoint: host,
			},
		},
	}

	rsp := &NewDIDResult{
		DID:             did,
		CreateOperation: signed_op,
		PrivateKey:      private_key,
	}

	return rsp, nil
}
