package decoder

import (
	"io"
)

type (
	// Decoder interface
	Decoder interface {
		Decode(
			body interface{},
			dest interface{},
		) error
		DecodeReader(
			reader io.Reader,
			dest interface{},
		) error
	}
)
