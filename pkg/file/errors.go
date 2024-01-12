package file

import "errors"

var (
	ErrInternal         = errors.New("server encountered an unexpected condition that prevented it from fulfilling the request")
	ErrFileNotFoundByID = errors.New("could not find any file with provided id")
)
