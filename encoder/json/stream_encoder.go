package json

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"

	gojsonencoder "github.com/ralvarezdev/go-json/encoder"
)

type (
	// StreamEncoder is the JSON encoder struct
	StreamEncoder struct{}
)

// NewStreamEncoder creates a new JSON encoder
//
// Returns:
//
//   - *StreamEncoder: The default encoder
func NewStreamEncoder() *StreamEncoder {
	return &StreamEncoder{}
}

// Encode encodes the body into JSON
//
// Parameters:
//
//   - body: The body to encode
//
// Returns:
//
//   - ([]byte): The encoded JSON
//   - error: The error if any
func (s StreamEncoder) Encode(
	body interface{},
) ([]byte, error) {
	// Check if body is nil
	if body == nil {
		return nil, gojsonencoder.ErrNilBody
	}

	// Create a buffer to write to
	buffer := new(bytes.Buffer)

	// Wrap it with a bufio.Writer
	writer := bufio.NewWriter(buffer)

	// Encode the body into JSON
	if err := json.NewEncoder(writer).Encode(body); err != nil {
		return nil, err
	}

	// Flush to ensure all data is written to the underlying buffer
	if err := writer.Flush(); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// EncodeAndWrite encodes the body into JSON and writes it to the writer
//
// Parameters:
//
//   - writer: The writer
//   - beforeWriteFn: The function to call before writing
//   - body: The body to encode
//
// Returns:
//
//   - error: The error if any
func (s StreamEncoder) EncodeAndWrite(
	writer io.Writer,
	beforeWriteFn func() error,
	body interface{},
) (err error) {
	// Check if the writer is nil
	if writer == nil {
		return gojsonencoder.ErrNilWriter
	}

	// Call the before write function if provided
	if beforeWriteFn != nil {
		if err = beforeWriteFn(); err != nil {
			return err
		}
	}

	// Encode the body into JSON
	if err = json.NewEncoder(writer).Encode(body); err != nil {
		return err
	}

	return nil
}
