package feed

import (
	"net/http"
)

func getTimelineHandler() (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

	}

	return http.HandlerFunc(fn), nil
}
