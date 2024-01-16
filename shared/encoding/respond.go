package encoding

import (
	"encoding/json"
	"errors"
	"net/http"
)

func Respond(w http.ResponseWriter, r *http.Request, data any, status int) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			return errors.New("failed to encode valid response")
		}
	}

	return nil
}

func ErrorRespond(w http.ResponseWriter, r *http.Request, status int, err error) {
	resp := NewErrorResponse(status, http.StatusText(status), err.Error())
	Respond(w, r, resp, status)
}
