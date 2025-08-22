package feed

import (
	"net/http"
)

func getFollowersHandler() (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

	}

	return http.HandlerFunc(fn), nil
}
