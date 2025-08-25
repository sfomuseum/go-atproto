package did

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/fxamacker/cbor/v2"
	"github.com/multiformats/go-multibase"
)

// https://atproto.com/specs/did
// https://atproto.com/specs/handle#handle-resolution

// For reference, this is what Blacksky does:
// https://github.com/blacksky-algorithms/rsky/blob/main/rsky-pds/src/plc/operations.rs#L50C5-L58C46
// https://github.com/blacksky-algorithms/rsky/blob/main/rsky-pds/src/plc/operations.rs#L281
// https://github.com/blacksky-algorithms/rsky/blob/main/rsky-common/src/sign.rs#L8

const DID_PLC string = "did:plc"
const DID_KEY string = "did:key"

const MB_ED25519 string = "\xED\x01"

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
	// RotationKeys        []string              `json:"rotationKeys"`
	AlsoKnownAs []string   `json:"alsoKnownAs"`
	Service     []*Service `json:"service"`
}

func (d *DID) String() string {
	return d.Id
}

func (d *DID) PublicKey(suffix string) ([]byte, error) {

	for _, m := range d.VerificationMethods {

		if !strings.HasSuffix(m.Id, suffix) {
			continue
		}

		key_mb := m.PublicKeyMultibase
		key_mb = strings.TrimLeft(key_mb, fmt.Sprintf("%s:", DID_KEY))

		_, body, err := multibase.Decode(key_mb)

		if err != nil {
			return nil, fmt.Errorf("Failed to decode multibase, %w", err)
		}

		return body[2:], nil
	}

	return nil, fmt.Errorf("Key not found")
}

func (d *DID) Marshal(wr io.Writer) error {
	enc := json.NewEncoder(wr)
	return enc.Encode(d)
}

type GenesisOperationService struct {
	Type     string `cbor:"type",json:"type"`
	Endpoint string `cbor:"endpoint,json:"endpoint""`
}

type GenesisOperation struct {
	Type                string                             `cbor:"type"`
	VerificationMethods map[string]string                  `cbor:"verificationMethods"`
	RotationKeys        []string                           `cbor:"rotationKeys"`
	AlsoKnownAs         []string                           `cbor:"alsoKnownAs"`
	Services            map[string]GenesisOperationService `cbor:"services"`
	Prev                interface{}                        `cbor:"prev"`
}

type GenesisOperationSigned struct {
	GenesisOperation
	Signature string `cbor:"sig"`
}

func NewDID(handle string, host string) (*DID, ed25519.PrivateKey, error) {

	// https://web.plc.directory/spec/v0.1/did-plc
	// In pseudo-code: did:plc:${base32Encode(sha256(createOp)).slice(0,24)}

	// Collect values for the essential operation data fields, including generating new secure key pairs if necessary

	public_key, private_key, err := ed25519.GenerateKey(rand.Reader)

	if err != nil {
		return nil, nil, fmt.Errorf("key generation: %w", err)
	}

	// Construct an “unsigned” regular operation object.
	// Include a prev field with null value. do not use the deprecated/legacy operation format for new DID creations

	// Base64url‑encoded public key – the spec uses this representation.
	public_b64 := base64.RawURLEncoding.EncodeToString(public_key)

	unsigned := GenesisOperation{
		Type: "plc_operation",
		VerificationMethods: map[string]string{
			"atproto": fmt.Sprintf("%s:%s", DID_KEY, public_b64),
		},
		RotationKeys: []string{
			fmt.Sprintf("%s:%s", DID_KEY, public_b64),
		},
		AlsoKnownAs: []string{fmt.Sprintf("at://%s", handle)},
		Services: map[string]GenesisOperationService{
			"atproto_pds": {
				Type:     "AtprotoPersonalDataServer",
				Endpoint: host,
			},
		},
		Prev: nil, // genesis – no previous operation
	}

	// Serialize the “unsigned” operation with DAG-CBOR, and sign the resulting bytes with one of the initial rotationKeys.
	// Encode the signature as base64url, and use that to construct a “signed” operation object

	enc_opts := cbor.CanonicalEncOptions()

	enc_mode, err := enc_opts.EncMode()

	if err != nil {
		return nil, nil, fmt.Errorf("cbor encoder: %w", err)
	}

	unsigned_b, err := enc_mode.Marshal(unsigned)

	if err != nil {
		return nil, nil, fmt.Errorf("unsigned CBOR marshal: %w", err)
	}

	sig := ed25519.Sign(private_key, unsigned_b)
	sig_b64 := base64.RawURLEncoding.EncodeToString(sig)

	// Serialize the “signed” operation with DAG-CBOR, take the SHA-256 hash of those bytes, and encode the hash bytes in base32.
	// use the first 24 characters to generate DID value (did:plc:<hashchars>)

	signed := GenesisOperationSigned{
		GenesisOperation: unsigned,
		Signature:        sig_b64,
	}

	signed_b, err := enc_mode.Marshal(signed)

	if err != nil {
		return nil, nil, fmt.Errorf("signed CBOR marshal: %w", err)
	}

	hash := sha256.Sum256(signed_b)

	b32_enc := base32.StdEncoding.WithPadding(base32.NoPadding)
	hash_b32 := b32_enc.EncodeToString(hash[:]) // 52 chars

	if len(hash_b32) < 24 {
		return nil, nil, fmt.Errorf("hash too short")
	}

	str_did := hash_b32[:24]

	// https://atproto.com/specs/cryptography

	combined := append([]byte(MB_ED25519), public_key...)

	public_mb, err := multibase.Encode(multibase.Base58BTC, combined)

	if err != nil {
		return nil, nil, err
	}

	// START OF code for sanity-checking the multibase encoding

	_, body, err := multibase.Decode(public_mb)

	if err != nil {
		return nil, nil, err
	}

	if len(body) < 2 {
		return nil, nil, fmt.Errorf("Decode key too short")
	}

	pk := ed25519.PublicKey(body[2:])

	if !pk.Equal(public_key) {
		return nil, nil, fmt.Errorf("Failed to encode/decode multibase public key")
	}

	// END OF code for sanity-checking the multibase encoding

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
		AlsoKnownAs: unsigned.AlsoKnownAs,
		Service: []*Service{
			&Service{
				Id:              "#atproto_pds",
				Type:            "AtprotoPersonalDataServer",
				ServiceEndpoint: host,
			},
		},
	}

	return did, private_key, nil
}
