package json

import (
	"errors"
)

var (
	ErrCodeNilDestination   string
	ErrCodeFailedToReadBody string
)

var (
	ErrNilDestination = errors.New("json destination is nil")
)
