package repo

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aaronland/go-http/v3/sanitize"
	"github.com/aaronland/go-http/v3/slog"
	"github.com/sfomuseum/go-pds"
)

const GetRecordHandlerURI string = "/xrpc/com.atproto.repo.getRecord"
const GetRecordHandlerMethod string = "GET"

type GetRecordResponse struct {
	CID       string           `json:"cid"`
	Record    *pds.Record      `json:"record"`
	BlockURI  string           `json:"blockUri"`
	Commit    string           `json:"commit"`
	Proof     *json.RawMessage `json:"proof,omitempty"` // optional, omitted here
	CreatedAt time.Time        `json:"createdAt"`
	UpdatedAt time.Time        `json:"updatedAt"`
}

type GetRecordHandlerOptions struct {
	RecordsDatabase pds.RecordsDatabase
}

func GetRecordHandler(opts *GetRecordHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		logger := slog.LoggerWithRequest(req, nil)

		if req.Method != GetRecordHandlerMethod {
			logger.Error("Method not allowed", "method", req.Method)
			http.Error(rsp, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// repo === did

		repo, err := sanitize.GetString(req, "repo")

		if err != nil {
			logger.Error("Invalid parameter", "parameter", "repo", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		if repo == "" {
			logger.Error("Missing parameter", "parameter", "repo")
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		logger = logger.With("repo", repo)

		collection, err := sanitize.GetString(req, "collection")

		if err != nil {
			logger.Error("Invalid parameter", "parameter", "collection", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		if collection == "" {
			logger.Error("Missing parameter", "parameter", "collection")
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		logger = logger.With("collection", collection)

		rkey, err := sanitize.GetString(req, "rkey")

		if err != nil {
			logger.Error("Invalid parameter", "parameter", "rkey", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		if rkey == "" {
			logger.Error("Missing parameter", "parameter", "rkey")
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		logger = logger.With("rkey", rkey)

		ctx := req.Context()

		rec, err := pds.GetRecord(ctx, opts.RecordsDatabase, repo, collection, rkey)

		if err != nil {
			logger.Error("Record not found", "error", err)
			http.Error(rsp, "Not found", http.StatusNotFound)
			return
		}

		get_rsp := GetRecordResponse{
			CID:       rec.CID,
			Record:    rec,
			BlockURI:  rec.BlockURI(),
			Commit:    rec.CID, // placeholder
			CreatedAt: time.Unix(rec.Created, 0),
			UpdatedAt: time.Unix(rec.LastModified, 0),
		}

		rsp.Header().Set("Content-type", "application/json")

		enc := json.NewEncoder(rsp)
		err = enc.Encode(get_rsp)

		if err != nil {
			logger.Error("Failed to encode record", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	return http.HandlerFunc(fn), nil
}
