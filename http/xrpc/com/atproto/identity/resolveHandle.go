package identity

import (
	"net/http"

	"github.com/aaronland/go-http/v3/sanitize"
	"github.com/aaronland/go-http/v3/slog"
	"github.com/sfomuseum/go-atproto"
	"github.com/sfomuseum/go-atproto/pds"
)

const ResolveHandleHandlerURI string = "/xrpc/com.atproto.identity.resolveHandle"
const ResolveHandleHandlerMethod string = "GET"

type ResolveHandleHandlerOptions struct {
	UsersDatabase pds.UsersDatabase
}

func ResolveHandleHandler(opts *ResolveHandleHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		logger := slog.LoggerWithRequest(req, nil)

		if req.Method != ResolveHandleHandlerMethod {
			logger.Error("Method not allowed", "method", req.Method)
			http.Error(rsp, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		handle, err := sanitize.GetString(req, "handle")

		if err != nil {
			logger.Error("Invalid parameter", "parameter", "handle", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		if handle == "" {
			logger.Error("Missing parameter", "parameter", "handle")
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		logger = logger.With("handle", handle)

		ctx := req.Context()

		rec, err := pds.GetUserWithHandle(ctx, opts.UsersDatabase, handle)

		if err != nil {

			if err == atproto.ErrNotFound {
				logger.Error("Record not found")
				http.Error(rsp, "Not found", http.StatusNotFound)
			} else {
				logger.Error("Failed to retrieve record", "error", err)
				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			}

			return
		}

		rsp.Write([]byte(rec.DID))
	}

	return http.HandlerFunc(fn), nil
}
