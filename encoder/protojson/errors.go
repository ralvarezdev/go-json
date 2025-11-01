package protojson

import (
	"errors"
)

const (
	ErrFieldNotHandled = "field not handled: %s"
	ErrFieldNotProtoMessage = "field is not a proto message: %s"
)

var (
	ErrNilStructInstance = errors.New("struct instance is nil")
)