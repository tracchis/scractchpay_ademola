package httputil

import (
	"encoding/json"
	"net/http"
	"strings"
)

// JSON render a generic interface as response of type json
func JSON(w http.ResponseWriter, code int, v interface{}) (err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	if v == nil || code == http.StatusNoContent {
		return
	}

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(true)

	if err = enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return
}

// JSONError render json error response
func JSONError(w http.ResponseWriter, status int, error string, messages map[string]string) {
	res := Response{
		Errors:   error,
		Messages: changeAttributeKeysInError(messages),
	}

	JSON(w, status, res)
	return
}

type AttributeErrors struct{}

// JSONSuccess render json success response
func JSONSuccess(w http.ResponseWriter, status int, responseBody interface{}) {
	JSON(w, status, responseBody)
	return
}

// Response is a generic response for APIs
type Response struct {
	Errors   string            `json:"error"`
	Messages map[string]string `json:"messages"`
}

func changeAttributeKeysInError(errorMessages map[string]string) map[string]string {
	newMap := make(map[string]string)

	for k, v := range errorMessages {
		keys := strings.Split(k, ".")
		key := keys[len(keys)-1]
		newMap[key] = v
		delete(errorMessages, k)
	}

	return newMap
}
