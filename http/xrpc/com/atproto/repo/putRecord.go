package repo

import (
	"net/http"
)

func putRecordHandler() (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

	}

	return http.HandlerFunc(fn), nil
}
