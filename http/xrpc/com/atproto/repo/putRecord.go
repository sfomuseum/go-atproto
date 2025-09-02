package repo

import (
	"net/http"

	"github.com/aaronland/go-http/v3/slog"
	"github.com/sfomuseum/go-atproto/pds"
)

const PutRecordHandlerURI string = "/xrpc/com.atproto.repo.putRecord"
const PutRecordHandlerMethod string = http.MethodPut

type PutRecordHandlerOptions struct {
	AccountsDatabase pds.AccountsDatabase
	RecordsDatabase  pds.RecordsDatabase
}

func PutRecordHandler() (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		logger := slog.LoggerWithRequest(req, nil)

		if req.Method != PutRecordHandlerMethod {
			logger.Error("Method not allowed", "method", req.Method)
			http.Error(rsp, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

	}

	return http.HandlerFunc(fn), nil
}
