package api

// https://web.plc.directory/api/redoc#operation/CreatePlcOp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/sfomuseum/go-atproto/plc"
)

type createPlcError struct {
	Message string `json:"message"`
}

func CreatePlc(ctx context.Context, str_did string, op plc.CreatePlcOperationSigned) error {

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

	rsp, err := http.DefaultClient.Do(req)

	if err != nil {
		return fmt.Errorf("Failed to execute request, %w", err)
	}

	defer rsp.Body.Close()

	if rsp.StatusCode != 200 {

		var m *createPlcError

		dec := json.NewDecoder(rsp.Body)
		err = dec.Decode(&m)

		if err != nil {
			slog.Error("Failed to decode response (error) body", "error", err)
			return fmt.Errorf("Request failed with error code %s", rsp.Status)
		}

		return fmt.Errorf("Request failed with error code %s: %s", rsp.Status, m.Message)
	}

	return nil
}
