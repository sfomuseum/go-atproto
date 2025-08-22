package repo

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aaronland/go-http/v3/sanitize"
	"github.com/sfomuseum/go-pds"
)

const GetRecordHandlerURI string = "/xrpc/com.atproto.repo.getRecord"

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

		if req.Method != http.MethodGet {
			http.Error(rsp, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// repo === did

		repo, err := sanitize.GetString(req, "repo")

		if err != nil {
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		if repo == "" {
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		collection, err := sanitize.GetString(req, "collection")

		if err != nil {
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		if collection == "" {
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		rkey, err := sanitize.GetString(req, "rkey")

		if err != nil {
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		if rkey == "" {
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		ctx := req.Context()

		rec, err := opts.RecordsDatabase.GetRecord(ctx, repo, collection, rkey)

		if err != nil {
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
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	return http.HandlerFunc(fn), nil
}
