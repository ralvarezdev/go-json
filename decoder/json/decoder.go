package json

import (
	"encoding/json"
	"io"

	gojsondecoder "github.com/ralvarezdev/go-json/decoder"
)

type (
	// Decoder struct
	Decoder struct{}
)

// NewDecoder creates a new JSON decoder
//
// Returns:
//
//   - *Decoder: The default decoder
func NewDecoder() *Decoder {
	return &Decoder{}
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
func (d Decoder) Decode(
	body any,
	dest any,
) error {
	// Check the body
	if body == nil {
		return gojsondecoder.ErrNilBody
	}

	// Check the body type
	reader, err := gojsondecoder.ToReader(body)
	if err != nil {
		return err
	}
	return d.DecodeReader(reader, dest)
}

// DecodeReader decodes the JSON body and stores it in the destination
//
// Parameters:
//
//   - reader: The reader to read the body from
//   - dest: The destination to store the decoded body
//
// Returns:
//
//   - error: The error if any
func (d Decoder) DecodeReader(
	reader io.Reader,
	dest any,
) error {
	// Check the reader
	if reader == nil {
		return gojsondecoder.ErrNilReader
	}

	// Check the decoder destination
	if dest == nil {
		return gojsondecoder.ErrNilDestination
	}

	// Get the body of the request
	body, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	// Decode JSON body into destination
	return json.Unmarshal(body, dest)
}
