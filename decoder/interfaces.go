package decoder

import (
	"io"
)

type (
	// Decoder interface
	Decoder interface {
		Decode(
			body any,
			dest any,
		) error
		DecodeReader(
			reader io.Reader,
			dest any,
		) error
	}
)
