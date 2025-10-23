package protojson

import (
	"encoding/json"
	"io"
	"reflect"

	goreflect "github.com/ralvarezdev/go-reflect"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// decodeNestedProtoMessages recursively decodes nested proto.Message fields in a struct
//
// Parameters:
//
// - body: The map representing the JSON data
// - dest: The destination struct to decode into
// - unmarshalOptions: Options for unmarshalling proto messages
//
// Returns:
//
// - error: The error if any
func decodeNestedProtoMessages(
	body map[string]any,
	dest any,
	unmarshalOptions *protojson.UnmarshalOptions,
) error {
	// Validate unmarshal options
	if unmarshalOptions == nil {
		unmarshalOptions = &protojson.UnmarshalOptions{}
	}

	// Dereference pointer if necessary
	v := reflect.ValueOf(dest)
	t := v.Type()
	if t.Kind() == reflect.Ptr {
		v = v.Elem()
		t = v.Type()
	}

	// Decode nested proto messages
	for i := 0; i < v.NumField(); i++ {
		// Get the field and its type
		field := v.Field(i)
		fieldType := t.Field(i)
		key := fieldType.Name

		// Check if the field can be set
		if !field.CanSet() {
			continue
		}

		// Check if the field is a proto.Message
		if protoField, ok := field.Interface().(proto.Message); ok {
			// Get the corresponding body field
			bodyField, exists := body[key]
			if !exists {
				continue
			}

			// Marshal the body field to JSON
			marshaledField, err := json.Marshal(bodyField)
			if err != nil {
				return err
			}

			// Unmarshal the JSON into the proto.Message field
			if unmarshalErr := unmarshalOptions.Unmarshal(
				marshaledField,
				protoField,
			); unmarshalErr != nil {
				return unmarshalErr
			}
		}

		// Check if the field is a struct
		if field.Kind() != reflect.Struct || field.Elem().Kind() != reflect.Struct {
			// Get the corresponding body field
			bodyField, exists := body[key]
			if !exists {
				continue
			}

			// Marshal the body field to JSON
			marshaledField, err := json.Marshal(bodyField)
			if err != nil {
				return err
			}

			// Unmarshal the JSON into the struct field
			if unmarshalErr := json.Unmarshal(
				marshaledField,
				field.Interface(),
			); unmarshalErr != nil {
				return unmarshalErr
			}
			continue
		}

		// Get the corresponding body field
		bodyField, exists := body[key]
		if !exists {
			continue
		}

		// Marshal the body field to JSON
		marshaledField, err := json.Marshal(bodyField)
		if err != nil {
			return err
		}

		// Initialize the map to hold nested body
		var nestedBody map[string]any
		if unmarshalErr := json.Unmarshal(marshaledField, &nestedBody); unmarshalErr != nil {
			return unmarshalErr
		}

		// Decode the nested body
		if decodeErr := decodeNestedProtoMessages(
			nestedBody,
			field.Addr().Interface(),
			unmarshalOptions,
		); decodeErr != nil {
			return decodeErr
		}
	}

	// Map the body to the struct field
	return goreflect.MapToStruct(body, dest)
}

// UnmarshalByReflection decodes JSON from io.Reader into a destination, handling nested proto.Message fields.
//
// Parameters:
//
//   - r: The io.Reader to read JSON data from
//   - dest: The destination to decode the JSON data into
//   - unmarshalOptions: Options for unmarshalling proto messages (optional, can be nil)
//
// Returns:
//
//   - error: The error if any
func UnmarshalByReflection(
	reader io.Reader,
	dest any,
	unmarshalOptions *protojson.UnmarshalOptions,
) error {
	// Validate unmarshal options
	if unmarshalOptions == nil {
		unmarshalOptions = &protojson.UnmarshalOptions{}
	}

	// Read all data from the reader
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	// Check if the destination is a proto.Message
	destProto, ok := dest.(proto.Message)
	if ok {
		// Directly unmarshal if it's a proto.Message
		return unmarshalOptions.Unmarshal(data, destProto)
	}

	// Initialize the map to hold intermediate JSON data
	var body map[string]any
	if unmarshalErr := json.Unmarshal(data, &body); unmarshalErr != nil {
		return unmarshalErr
	}

	return decodeNestedProtoMessages(body, dest, unmarshalOptions)
}
