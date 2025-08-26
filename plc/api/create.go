package api

// https://web.plc.directory/api/redoc#operation/CreatePlcOp
// https://github.com/did-method-plc/did-method-plc/blob/main/packages/server/src/routes.ts#L114
// https://github.com/did-method-plc/did-method-plc/blob/main/packages/server/src/constraints.ts#L21 	<-- validateIncomingOp
// https://github.com/did-method-plc/did-method-plc/blob/main/packages/server/src/db/index.ts#L101 	<-- validateAndAddOp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	_ "log/slog"
	"os"

	"github.com/sfomuseum/go-atproto/plc"
)

type createPlcError struct {
	Message string `json:"message"`
}

// https://github.com/bluesky-social/indigo/blob/main/plc/client.go#L61
// https://github.com/did-method-plc/did-method-plc/blob/944a9ca36dd06b11630ec4d069c1b70fc6961ccf/website/spec/plc-server-openapi3.yaml#L244

func Create(ctx context.Context, str_did string, op plc.PlcOperationSigned) error {

	enc_op, err := json.Marshal(op)

	if err != nil {
		return err
	}

	op_r := bytes.NewReader(enc_op)

	u := NewURL()
	u.Path = fmt.Sprintf("/%s", str_did)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), op_r)

	if err != nil {
		return fmt.Errorf("Failed to create new request, %w", err)
	}

	req.Header.Set("Content-type", "application/json")

	// START OF DEBUGGING

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(op)

	if err != nil {
		return err
	}

	req2 := req.Clone(ctx)
	req2.Write(os.Stdout)

	os.Stdout.Write([]byte("\n\n"))

	// END OF DEBUGGING

	// Currently failing here because... ???
	// https://github.com/did-method-plc/did-method-plc/blob/main/packages/server/src/constraints.ts#L30-L45

	rsp, err := http.DefaultClient.Do(req)

	if err != nil {
		return fmt.Errorf("Failed to execute request, %w", err)
	}

	defer rsp.Body.Close()

	if rsp.StatusCode != 200 {

		/*
			var m *createPlcError

			body, _ := io.ReadAll(rsp.Body)

			err := json.Unmarshal(body, &m)

			if err != nil {
				slog.Error("Failed to decode response (error) body", "error", err)
				return fmt.Errorf("Request failed with error code %s", rsp.Status)
			}
		*/

		return fmt.Errorf("Request failed with error code %d: %s", rsp.StatusCode, rsp.Status)
	}

	return nil
}
