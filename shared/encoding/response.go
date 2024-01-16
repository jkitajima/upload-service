package encoding

import "strings"

type DataResponse struct {
	Data any `json:"data"`
}

type ErrorResponse struct {
	Error *ErrorObject `json:"error"`
}

func NewErrorResponse(status int, title, desc string) ErrorResponse {
	if desc != "" {
		desc = strings.ToUpper(desc[:1]) + desc[1:] + "."
	}
	return ErrorResponse{&ErrorObject{status, title, desc}}
}

type ErrorObject struct {
	Status      int    `json:"status"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
