package protojson

import (
	"io"

	gojsondecoder "github.com/ralvarezdev/go-json/decoder"
	"google.golang.org/protobuf/encoding/protojson"
)

type (
	Decoder struct {
		unmarshalOptions protojson.UnmarshalOptions
	}
)

// NewDecoder creates a new Decoder instance
//
// Returns:
//
//   - *Decoder: The decoder instance
func NewDecoder() *Decoder {
	// Initialize unmarshal options
	unmarshalOptions := protojson.UnmarshalOptions{
		DiscardUnknown: true,
		AllowPartial:   true,
	}

	return &Decoder{
		unmarshalOptions: unmarshalOptions,
	}
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
	body interface{},
	dest interface{},
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

// DecodeReader  decodes a JSON body from a reader into a destination
//
// Parameters:
//
//   - reader: The io.Reader to read the body from
//   - dest: The destination to decode the body into
//
// Returns:
//
//   - error: The error if any
func (d Decoder) DecodeReader(
	reader io.Reader,
	dest interface{},
) error {
	// Check the reader
	if reader == nil {
		return gojsondecoder.ErrNilReader
	}

	// Check the decoder destination
	if dest == nil {
		return gojsondecoder.ErrNilDestination
	}

	return UnmarshalByReflection(
		reader,
		dest,
		&d.unmarshalOptions,
	)
}
