package decoder

import (
	"errors"
)

var (
	ErrInvalidInstance = errors.New("invalid instance provided to create a reader")
	ErrNilBody         = errors.New("body cannot be nil")
	ErrNilReader       = errors.New("reader cannot be nil")
	ErrNilDestination  = errors.New("destination cannot be nil")
	ErrNilDecoder      = errors.New("decoder is nil")
)
