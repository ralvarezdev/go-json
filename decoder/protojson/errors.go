package protojson

import (
	"errors"
)

const (
	ErrFieldNotHandled      = "field not handled on decoding: %s"
	ErrFieldNotProtoMessage = "field %s is not a proto message"
)

var (
	ErrDestinationNotProtoMessage = errors.New("destination is not a proto message")
	ErrNilReader                  = errors.New("nil reader")
	ErrNilMapper                  = errors.New("decoder mapper is nil")
	ErrNilDestinationInstance     = errors.New("nil destination instance")
	ErrNilDestination             = errors.New("nil destination")
)
