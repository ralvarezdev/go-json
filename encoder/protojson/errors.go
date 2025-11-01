package protojson

import (
	"errors"
)

const (
	ErrFieldNotHandled = "field not handled on encoding: %s"
	ErrFieldNotProtoMessage = "field is not a proto message: %s"
)

var (
	ErrNilBody  = errors.New("body is nil")
	ErrNilMapper = errors.New("encoder mapper is nil")
	ErrNilStructInstance = errors.New("struct instance is nil")
)