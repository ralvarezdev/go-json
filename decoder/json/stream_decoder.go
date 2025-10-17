package json

import (
	"encoding/json"
	"io"

	gojsondecoder "github.com/ralvarezdev/go-json/decoder"
)

type (
	// StreamDecoder is the JSON decoder struct
	StreamDecoder struct{}
)

// NewStreamDecoder creates a new JSON decoder
//
// Returns:
//
//   - *StreamDecoder: The default decoder
func NewStreamDecoder() *StreamDecoder {
	return &StreamDecoder{}
}

// Decode decodes the JSON body from an any value and stores it in the destination
//
// Parameters:
//
//   - body: The body to decode
//   - dest: The destination to store the decoded body
//
// Returns:
//
//   - error: The error if any
func (s StreamDecoder) Decode(
	body interface{},
	dest interface{},
) error {
	// Check the body type
	reader, err := gojsondecoder.ToReader(body)
	if err != nil {
		return err
	}
	return s.DecodeReader(reader, dest)
}

// DecodeReader decodes a JSON body from a reader into a destination
//
// Parameters:
//
//   - reader: The reader to read the body from
//   - dest: The destination to store the decoded body
//
// Returns:
//
//   - error: The error if any
func (s StreamDecoder) DecodeReader(
	reader io.Reader,
	dest interface{},
) error {
	// Check the decoder destination
	if dest == nil {
		return ErrNilDestination
	}

	// Create the stream decoder
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()

	// Decode JSON body into destination
	return decoder.Decode(dest)
}
