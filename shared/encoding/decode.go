package encoding

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("encoding: %w", err)
	}
	return v, nil
}
