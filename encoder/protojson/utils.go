package protojson

import (
	"encoding/json"
	"reflect"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// processField processes a single struct field for marshaling
//
// Parameters:
//
//   - field: The reflect.Value of the field
//   - fieldType: The reflect.StructField of the field
//   - marshalOptions: The protojson.MarshalOptions to use
//   - result: The map to store the processed field
//
// Returns:
//
// - error: The error if any
func processField(
	field reflect.Value,
	fieldType *reflect.StructField,
	marshalOptions *protojson.MarshalOptions,
	result map[string]any,
) error {
	// Check if fieldType is nil
	if fieldType == nil {
		return nil
	}

	// Get the field value
	val := field.Interface()

	// Handle proto.Message fields
	switch fv := val.(type) {
	case proto.Message:
		// Marshal proto.Message to JSON
		data, err := marshalOptions.Marshal(fv)
		if err != nil {
			return err
		}

		// Store as json.RawMessage to avoid double encoding
		result[fieldType.Name] = json.RawMessage(data)
	default:
		// Recursively handle nested structs
		if field.Kind() == reflect.Struct {
			nested, err := PrecomputeMarshalByReflection(
				field,
				marshalOptions,
			)
			if err != nil {
				return err
			}
			result[fieldType.Name] = nested
		} else {
			result[fieldType.Name] = val
		}
	}
	return nil
}

// PrecomputeMarshalByReflection marshals a struct to a map[string]any using reflection, handling nested
// proto.Message fields appropriately
//
// Parameters:
//
//   - v: The reflect.Value of the struct to marshal
//   - marshalOptions: The protojson.MarshalOptions to use (optional)
//
// Returns:
//
//   - map[string]any: The marshaled struct as a map
//   - error: The error if any
func PrecomputeMarshalByReflection(
	v reflect.Value,
	marshalOptions *protojson.MarshalOptions,
) (map[string]any, error) {
	// Ensure marshalOptions is not nil
	if marshalOptions == nil {
		marshalOptions = &protojson.MarshalOptions{}
	}

	// Dereference pointer if necessary
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Prepare result map
	result := make(map[string]any)
	t := v.Type()

	// Handle nested proto.Message fields
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Check if the field can be interfaced
		if !field.CanInterface() {
			continue
		}

		// Process the field
		if err := processField(
			field,
			&fieldType,
			marshalOptions,
			result,
		); err != nil {
			return nil, err
		}
	}
	return result, nil
}
