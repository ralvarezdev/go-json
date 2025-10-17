package protojson

import (
	"io"
	"reflect"

	gojsonencoder "github.com/ralvarezdev/go-json/encoder"
	gojsonencoderjson "github.com/ralvarezdev/go-json/encoder/json"
	"google.golang.org/protobuf/encoding/protojson"
)

type (
	// Encoder is the implementation of the Encoder interface
	Encoder struct {
		jsonEncoder    *gojsonencoderjson.Encoder
		marshalOptions protojson.MarshalOptions
	}
)

// NewEncoder creates a new Encoder instance
//
// Returns:
//
// - *Encoder: the new Encoder instance
func NewEncoder() *Encoder {
	// Initialize the JSON encoder
	jsonEncoder := gojsonencoderjson.NewEncoder()

	// Initialize unmarshal options
	marshalOptions := protojson.MarshalOptions{
		AllowPartial: true,
	}

	return &Encoder{
		jsonEncoder:    jsonEncoder,
		marshalOptions: marshalOptions,
	}
}

// PrecomputeMarshal precomputes the marshaled body by reflecting on the instance
//
// Parameters:
//
// - body: The body to precompute the marshaled body for
//
// Returns:
//
// - (map[string]interface{}, error): The precomputed marshaled body and the error if any
func (e Encoder) PrecomputeMarshal(
	body interface{},
) (map[string]interface{}, error) {
	// Reflect on the instance to get its fields
	v := reflect.ValueOf(body)

	// Precompute the marshaled body
	precomputedMarshal, err := PrecomputeMarshalByReflection(
		v,
		&e.marshalOptions,
	)
	if err != nil {
		return nil, err
	}
	return precomputedMarshal, nil
}

// Encode encodes the given body to JSON
//
// Parameters:
//
//   - body: The body to encode
//
// Returns:
//
//   - ([]byte, error): The encoded body and the error if any
func (e Encoder) Encode(
	body interface{},
) ([]byte, error) {
	// Check if body is nil
	if body == nil {
		return nil, gojsonencoder.ErrNilBody
	}

	// Marshal the instance to get the precomputed body
	precomputedMarshal, err := e.PrecomputeMarshal(body)
	if err != nil {
		return nil, err
	}
	return e.jsonEncoder.Encode(precomputedMarshal)
}

// EncodeAndWrite encodes and writes the given body to the writer
//
// Parameters:
//
//   - writer: The writer to write the encoded body to
//   - beforeWriteFn: The function to call before writing the body
//   - body: The body to encode
//
// Returns:
//
// - error: The error if any
func (e Encoder) EncodeAndWrite(
	writer io.Writer,
	beforeWriteFn func() error,
	body interface{},
) error {
	// Check if the writer is nil
	if writer == nil {
		return gojsonencoder.ErrNilWriter
	}

	// Marshal the instance to get the precomputed body
	precomputedMarshal, err := e.PrecomputeMarshal(body)
	if err != nil {
		return err
	}
	return e.jsonEncoder.EncodeAndWrite(
		writer,
		beforeWriteFn,
		precomputedMarshal,
	)
}
