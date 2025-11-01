package protojson

import (
	"io"

	"google.golang.org/protobuf/encoding/protojson"

	gojsondecoder "github.com/ralvarezdev/go-json/decoder"
)

type (
	Decoder struct {
		unmarshalOptions protojson.UnmarshalOptions
		cache bool
		cachedMappers map[string]*Mapper
	}
	
	// Options are the additional settings for the decoder implementation
	Options struct {
		// cache indicates whether to cache the precompute unmarshal by reflection functions
		cache bool
	}
)

// NewOptions creates a new Options instance
// 
// Parameters:
// 
//  - cache: indicates whether to cache the precompute unmarshal by reflection functions
// 
// Returns:
// 
// - *Options: the new Options instance
func NewOptions(
	cache bool,
) *Options {
	return &Options{
		cache: cache,
	}
}

// NewDecoder creates a new Decoder instance
// 
// Parameters:
// 
//  - options: the additional settings for the decoder implementation
//
// Returns:
//
//   - *Decoder: The decoder instance
func NewDecoder(options *Options) *Decoder {
	// Initialize cache setting
	cache := false
	if options != nil {
		cache = options.cache
	}
	
	// Initialize unmarshal options
	unmarshalOptions := protojson.UnmarshalOptions{
		DiscardUnknown: true,
		AllowPartial:   true,
	}

	return &Decoder{
		unmarshalOptions: unmarshalOptions,
		cache:            cache,
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
	
	// Read all data from the reader
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	return UnmarshalByReflection(
		data,
		dest,
		&d.unmarshalOptions,
	)
}
