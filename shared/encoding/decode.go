package encoding

import (
	"encoding/json"
	"net/http"
)

func Decode(w http.ResponseWriter, r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}
