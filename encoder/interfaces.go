package encoder

import (
	"io"
)

type (
	// Encoder interface
	Encoder interface {
		Encode(
			body any,
		) ([]byte, error)
		EncodeAndWrite(
			writer io.Writer,
			beforeWriteFn func() error,
			body any,
		) error
	}

	// ProtoJSONEncoder interface
	ProtoJSONEncoder interface {
		Encoder
		PrecomputeMarshal(
			body any,
		) (map[string]any, error)
	}
)
