package httphandler

import "upload/pkg/file"

type DataResponse struct {
	Data *file.File `json:"data"`
}

type ErrorsResponse struct {
	Errors []*ErrorObject `json:"errors"`
}

type ErrorObject struct {
	Code        uint16 `json:"code"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func NewErrorsResponse(errors ...*ErrorObject) ErrorsResponse {
	err := make([]*ErrorObject, 0, len(errors))
	err = append(err, errors...)
	return ErrorsResponse{err}
}
