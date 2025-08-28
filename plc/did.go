package plc

import (
	"context"
	"fmt"
	"strings"

	"github.com/bluesky-social/indigo/atproto/crypto"
	_ "github.com/bluesky-social/indigo/atproto/data"
	"github.com/bluesky-social/indigo/atproto/identity"
	"github.com/bluesky-social/indigo/atproto/syntax"
	"github.com/did-method-plc/go-didplc"
	"github.com/whyrusleeping/go-did"
)

// https://atproto.com/specs/did

// For reference, this is what Blacksky does:
// https://github.com/blacksky-algorithms/rsky/blob/main/rsky-pds/src/plc/operations.rs#L50C5-L58C46
// https://github.com/blacksky-algorithms/rsky/blob/main/rsky-pds/src/plc/operations.rs#L281
// https://github.com/blacksky-algorithms/rsky/blob/main/rsky-common/src/sign.rs#L8

type NewDIDResult struct {
	DID          *identity.DIDDocument
	PlcOperation didplc.RegularOp
	PrivateKey   *crypto.PrivateKeyK256
}

func NewDID(ctx context.Context, service string, handle string) (*NewDIDResult, error) {

	handle = strings.TrimPrefix(handle, "at://")
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
		fmt.Sprintf("at://%s", parsed_handle),
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
		return nil, fmt.Errorf("Failed to verify signature for op, %w", err)
	}

	did_id, err := op.DID()

	if err != nil {
		return nil, fmt.Errorf("Failed to derive DID for op, %w", err)
	}

	_, err = did.ParseDID(did_id)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse plc did, %w", err)
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

	rsp := &NewDIDResult{
		DID:          doc,
		PlcOperation: op,
		PrivateKey:   private_key,
	}

	return rsp, nil
}
