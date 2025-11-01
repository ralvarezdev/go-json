package protojson

import (
	"io"

	"google.golang.org/protobuf/encoding/protojson"

	goreflect "github.com/ralvarezdev/go-reflect"

	gojsondecoder "github.com/ralvarezdev/go-json/decoder"
)

type (
	Decoder struct {
		unmarshalOptions protojson.UnmarshalOptions
		cache            bool
		cachedMappers    map[string]*Mapper
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
//   - cache: indicates whether to cache the precompute unmarshal by reflection functions
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
//   - options: the additional settings for the decoder implementation
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
//   - reader: The body to decode
//   - dest: The destination to store the decoded body
//
// Returns:
//
//   - error: The error if any
func (d Decoder) Decode(
	reader any,
	dest any,
) error {
	// Check the reader
	if reader == nil {
		return ErrNilReader
	}

	// Check the reader type
	parsedReader, err := gojsondecoder.ToReader(reader)
	if err != nil {
		return err
	}

	return d.DecodeReader(parsedReader, dest)
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

	// Read all body from the reader
	body, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	// Check if the cache is enabled and use cached mapper if available
	if d.cache && d.cachedMappers != nil {
		// Get the unique type identifier for the destination
		uniqueTypeReference := goreflect.UniqueTypeReference(dest)

		// Check if there is a cached mapper for the destination type
		if mapper, found := d.cachedMappers[uniqueTypeReference]; found {
			return mapper.UnmarshalByReflection(
				body,
				dest,
				&d.unmarshalOptions,
			)
		}
	}

	// Initialize the cache map if caching is enabled
	if d.cache && d.cachedMappers == nil {
		d.cachedMappers = make(map[string]*Mapper)
	}

	// Create a new mapper for the destination type
	mapper, err := NewMapper(dest)
	if err != nil {
		return err
	}

	// Store the mapper in the cache if caching is enabled
	if d.cache {
		uniqueTypeReference := goreflect.UniqueTypeReference(dest)
		d.cachedMappers[uniqueTypeReference] = mapper
	}

	// Unmarshal the body into the destination using the mapper
	return mapper.UnmarshalByReflection(
		body,
		dest,
		&d.unmarshalOptions,
	)
}
