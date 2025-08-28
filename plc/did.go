package plc

import (
	"context"
	"fmt"
	"strings"

	"github.com/bluesky-social/indigo/atproto/crypto"
	"github.com/bluesky-social/indigo/atproto/identity"
	"github.com/bluesky-social/indigo/atproto/syntax"
	"github.com/did-method-plc/go-didplc"
)

// NewDIDResults is a struct wrapping items created by the `NewDID` method.
type NewDIDResult struct {
	// The `identity.DIDDocument` instance wrapping a given handle at a service.
	DID *identity.DIDDocument
	// The signed PLC (regular) operation used to create the DID which can be submitted to a PLC directory service as a separate task.
	Operation didplc.Operation
	// The private signing key that was created for the new DID.
	PrivateKey *crypto.PrivateKeyK256
}

// NewDID generates a new `identity.DIDDocument` for 'handle' at 'service' and returns a signed `didplc.Operation`
// which can be submitted to a PLC directory service as a separate task. The identity document, signed operations
// as well as the private signing key associated with the DID are returned in a `NewDIDResult` struct.
func NewDID(ctx context.Context, service string, handle string) (*NewDIDResult, error) {

	// This basically follows the same logic/code defined in bluesky-social/goat
	// https://github.com/bluesky-social/goat/blob/main/plc.go#L416

	handle = strings.TrimPrefix(handle, AT_SCHEME)
	parsed_handle, err := syntax.ParseHandle(handle)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse handler, %w", err)
	}

	parsed_handle = parsed_handle.Normalize()

	// https://web.plc.directory/spec/v0.1/did-plc
	// In pseudo-code: did:plc:${base32Encode(sha256(createOp)).slice(0,24)}

	// Collect values for the essential operation data fields, including generating new secure key pairs if necessary
	// Only secp256k1 (“k256”) and NIST P-256 (“p256”) keys are currently supported for rotation keys, whereas verificationMethods keys can be any syntactically-valid did:key.

	private_key, err := crypto.GeneratePrivateKeyK256()

	if err != nil {
		return nil, fmt.Errorf("Failed to generate private key, %w", err)
	}

	public_key, err := private_key.PublicKey()

	if err != nil {
		return nil, fmt.Errorf("Failed to derive public key, %w", err)
	}

	// Construct an “unsigned” regular operation object.
	// Include a prev field with null value. do not use the deprecated/legacy operation format for new DID creations

	services := map[string]didplc.OpService{
		"atproto_pds": didplc.OpService{
			Type:     "AtprotoPersonalDataServer",
			Endpoint: service,
		},
	}

	verification_methods := map[string]string{
		"atproto": public_key.DIDKey(),
	}

	also_known_as := []string{
		fmt.Sprintf("%s%s", AT_SCHEME, parsed_handle),
	}

	rotation_keys := []string{
		public_key.DIDKey(),
	}

	op := didplc.RegularOp{
		Type:                "plc_operation",
		RotationKeys:        rotation_keys,
		VerificationMethods: verification_methods,
		AlsoKnownAs:         also_known_as,
		Services:            services,
	}

	err = op.Sign(private_key)

	if err != nil {
		return nil, fmt.Errorf("Failed to sign operation, %w", err)
	}

	err = op.VerifySignature(public_key)

	if err != nil {
		return nil, fmt.Errorf("Failed to verify signature for operation, %w", err)
	}

	did_id, err := op.DID()

	if err != nil {
		return nil, fmt.Errorf("Failed to derive DID from operation, %w", err)
	}

	doc := &identity.DIDDocument{
		DID:         syntax.DID(did_id),
		AlsoKnownAs: op.AlsoKnownAs,
		VerificationMethod: []identity.DocVerificationMethod{
			identity.DocVerificationMethod{
				ID:                 fmt.Sprintf("%s#atproto", did_id),
				Type:               "Multikey",
				Controller:         did_id,
				PublicKeyMultibase: public_key.DIDKey(),
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

	oe := didplc.OpEnum{
		Regular: &op,
	}

	as_op := oe.AsOperation()

	if as_op == nil {
		return nil, fmt.Errorf("Failed to derive as operation")
	}

	rsp := &NewDIDResult{
		DID:        doc,
		Operation:  as_op,
		PrivateKey: private_key,
	}

	return rsp, nil
}

func TombstoneDID(ctx context.Context, doc *identity.DIDDocument, prev string, private_key *crypto.PrivateKeyK256) error {

	op := didplc.TombstoneOp{
		Type: "plc_tombstone",
		Prev: prev,
	}

	err := op.Sign(private_key)

	if err != nil {
		return fmt.Errorf("Failed to sign op, %w", err)
	}

	var atproto_key string

	for _, m := range doc.VerificationMethod {

		if strings.HasSuffix(m.ID, "#atproto") && m.Controller == doc.DID.String() {
			atproto_key = m.PublicKeyMultibase
			break
		}
	}

	if atproto_key == "" {
		return fmt.Errorf("Missing atproto verification method")
	}

	public_key, err := crypto.ParsePublicMultibase(atproto_key)

	if err != nil {
		return fmt.Errorf("Failed to derive public key, %w", err)
	}

	err = op.VerifySignature(public_key)

	if err != nil {
		return fmt.Errorf("Failed to verify signature for operation, %w", err)
	}

	return nil
}
