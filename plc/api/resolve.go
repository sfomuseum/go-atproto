package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bluesky-social/indigo/atproto/identity"
)

// https://web.plc.directory/api/redoc#operation/ResolveDid

func ResolveDID(ctx context.Context, str_did string) (*identity.DIDDocument, error) {

	u := NewURL()
	u.Path = fmt.Sprintf("/%s", str_did)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new request, %w", err)
	}

	rsp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("Failed to execute request, %w", err)
	}

	defer rsp.Body.Close()

	if rsp.StatusCode != 200 {
		return nil, fmt.Errorf("Request failed with error code %d %s", rsp.StatusCode, rsp.Status)
	}

	var doc *identity.DIDDocument

	dec := json.NewDecoder(rsp.Body)
	err = dec.Decode(&doc)

	if err != nil {
		return nil, fmt.Errorf("Failed to decode DID, %w", err)
	}

	return doc, nil
}
