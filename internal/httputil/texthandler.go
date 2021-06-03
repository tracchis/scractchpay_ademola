package httputil

import (
	"fmt"
	"net/http"
)

// TextHandler returns a HandlerFunc that writes a constant Content-Type and string as a response.
func TextHandler(status int, contentType, response string) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		if contentType != "" {
			wr.Header().Set("Content-Type", contentType)
		}

		wr.WriteHeader(status)
		fmt.Fprint(wr, response)
	}
}
