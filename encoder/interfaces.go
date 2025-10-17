package encoder

import (
	"io"
)

type (
	// Encoder interface
	Encoder interface {
		Encode(
			body interface{},
		) ([]byte, error)
		EncodeAndWrite(
			writer io.Writer,
			beforeWriteFn func() error,
			body interface{},
		) error
	}

	// ProtoJSONEncoder interface
	ProtoJSONEncoder interface {
		PrecomputeMarshal(
			body interface{},
		) (map[string]interface{}, error)
	}
)
