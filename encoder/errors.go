package encoder

import (
	"errors"
)

var (
	ErrNilWriter  = errors.New("writer is nil")
	ErrNilBody    = errors.New("body cannot be nil")
	ErrNilEncoder = errors.New("encoder is nil")
)
