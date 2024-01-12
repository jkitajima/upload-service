package blob

import "errors"

var (
	ErrInternal = errors.New("error while communicating with blob storage")
	ErrNotFound = errors.New("blob was not found")
)
