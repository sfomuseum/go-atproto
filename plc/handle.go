package plc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type resolveHandleResponse struct {
	DID string `json:"did"`
}

// ResolveHandle resolve a handle (composed of 'handle' + "." + 'service') to its unique DID identifer by
// querying the "com.atproto.identity.resolveHandle" endpoint of 'service'.
func ResolveHandle(ctx context.Context, service string, handle string) (string, error) {

	q := url.Values{}
	q.Set("handle", handle)

	u, err := url.Parse(service)

	if err != nil {
		return "", fmt.Errorf("Failed to parse service URL, %w", err)
	}

	u.Path = "/xrpc/com.atproto.identity.resolveHandle"
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)

	if err != nil {
		return "", fmt.Errorf("Failed to create request, %w", err)
	}

	req.Header.Set("Content-type", "application/json")

	rsp, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", fmt.Errorf("Failed to execute request, %w", err)
	}

	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Query failed with code %d %s", rsp.StatusCode, rsp.Status)
	}

	var did_rsp *resolveHandleResponse

	dec := json.NewDecoder(rsp.Body)
	err = dec.Decode(&did_rsp)

	if err != nil {
		return "", fmt.Errorf("Failed to decode response, %w", err)
	}

	return did_rsp.DID, nil
}
