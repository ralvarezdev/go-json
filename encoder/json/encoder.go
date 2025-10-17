package json

import (
	"encoding/json"
	"io"
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
	body interface{},
) ([]byte, error) {
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
	body interface{},
) error {
	// Encode the body into JSON
	jsonBody, err := e.Encode(body)
	if err != nil {
		return err
	}

	// Call the before write function if provided
	if beforeWriteFn != nil {
		if err = beforeWriteFn(); err != nil {
			return err
		}
	}

	// Write the JSON body to the writer
	_, err = writer.Write(jsonBody)
	return err
}
