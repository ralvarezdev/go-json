package json

import (
	"encoding/json"
	"io"

	gojsonencoder "github.com/ralvarezdev/go-json/encoder"
)

type (
	// Encoder struct
	Encoder struct{}
)

// NewEncoder creates a new default JSON encoder
//
// Returns:
//
//   - *Encoder: The default encoder
func NewEncoder() *Encoder {
	return &Encoder{}
}

// Encode encodes the body into JSON bytes
//
// Parameters:
//
//   - body: The body to encode
//
// Returns:
//
//   - []byte: The encoded JSON bytes
//   - error: The error if any
func (e Encoder) Encode(
	body any,
) ([]byte, error) {
	// Check if body is nil
	if body == nil {
		return nil, gojsonencoder.ErrNilBody
	}

	// Marshal the body into JSON
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return jsonBody, nil
}

// EncodeAndWrite encodes the body and writes it to the writer
//
// Parameters:
//
// - writer: The writer to write the response to
// - beforeWriteFn: The function to call before writing the response
// - body: The body to encode
//
// Returns:
//
// - error: The error if any
func (e Encoder) EncodeAndWrite(
	writer io.Writer,
	beforeWriteFn func() error,
	body any,
) error {
	// Check if the writer is nil
	if writer == nil {
		return gojsonencoder.ErrNilWriter
	}

	// Encode the body into JSON
	jsonBody, err := e.Encode(body)
	if err != nil {
		return err
	}

	// Call the before write function if provided
	if beforeWriteFn != nil {
		if fnErr := beforeWriteFn(); fnErr != nil {
			return fnErr
		}
	}

	// Write the JSON body to the writer
	_, writeErr := writer.Write(jsonBody)
	return writeErr
}
