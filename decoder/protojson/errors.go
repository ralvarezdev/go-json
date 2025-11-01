package protojson

import (
	"errors"
)

const (
	ErrFieldNotHandled = "field not handled on decoding: %s"
)

var (
	ErrNilReader 			= errors.New("nil reader")
	ErrNilMapper 			= errors.New("decoder mapper is nil")
	ErrNilDestinationInstance = errors.New("nil destination instance")
	ErrNilDestination = 		errors.New("nil destination")
)