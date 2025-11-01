package protojson

import (
	"io"

	"google.golang.org/protobuf/encoding/protojson"

	gojsonencoder "github.com/ralvarezdev/go-json/encoder"
	gojsonencoderjson "github.com/ralvarezdev/go-json/encoder/json"
	goreflect "github.com/ralvarezdev/go-reflect"
)

type (
	// Encoder is the implementation of the Encoder interface
	Encoder struct {
		jsonEncoder    *gojsonencoderjson.Encoder
		marshalOptions protojson.MarshalOptions
		cache bool
		cachedMappers map[string]*Mapper
	}
	
	// Options are the additional settings for the encoder implementation
	Options struct {
		// cache indicates whether to cache the precompute marshal by reflection functions
		cache bool
	}
)

// NewOptions creates a new Options instance
//
// Parameters:
// 
// - cache: indicates whether to cache the precompute marshal by reflection functions
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

// NewEncoder creates a new Encoder instance
//
// Parameters:
// 
// - options: the additional settings for the encoder implementation
// 
// Returns:
//
// - *Encoder: the new Encoder instance
func NewEncoder(options *Options) *Encoder {
	// Initialize cache setting
	cache := false
	if options != nil {
		cache = options.cache
	}
	
	// Initialize the JSON encoder
	jsonEncoder := gojsonencoderjson.NewEncoder()

	// Initialize unmarshal options
	marshalOptions := protojson.MarshalOptions{
		AllowPartial: true,
	}

	return &Encoder{
		jsonEncoder:    jsonEncoder,
		marshalOptions: marshalOptions,
		cache: cache,
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
// - (map[string]any, error): The precomputed marshaled body and the error if any
func (e Encoder) PrecomputeMarshal(
	body any,
) (map[string]any, error) {
	// Check if body is nil
	if body == nil {
		return nil, gojsonencoder.ErrNilBody
	}
	
	// Check if the cache is true, if so try to get the mapper from the cache
	if e.cache && e.cachedMappers != nil {
		// Get the unique type identifier for the body
		uniqueTypeReference := goreflect.UniqueTypeReference(body)
		
		// Check if the mapper exists in the cache
		if mapper, ok := e.cachedMappers[uniqueTypeReference]; ok {
			precomputedMarshal, err := mapper.PrecomputeMarshalByReflection(body, &e.marshalOptions)
			if err != nil {
				return nil, err
			}
			return precomputedMarshal, nil
		}
	}
	
	// Create the cache map if it doesn't exist
	if e.cache && e.cachedMappers == nil {
		e.cachedMappers = make(map[string]*Mapper)
	}
	
	// Create a new mapper and store it in the cache if caching is enabled
	mapper, err := NewMapper(body)
	if err != nil {
		return nil, err
	}
	if e.cache {
		// Get the unique type identifier for the body
		uniqueTypeReference := goreflect.UniqueTypeReference(body)
		
		// Store the mapper in the cache
		e.cachedMappers[uniqueTypeReference] = mapper
	}

	// Marshal the instance to get the precomputed body
	precomputedMarshal, err := mapper.PrecomputeMarshalByReflection(
		body,
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
	body any,
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
	body any,
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
