package protojson

import (
	"encoding/json"
	"reflect"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// PrecomputeMarshalByReflection marshals a struct to a map[string]interface{} using reflection, handling nested proto.Message fields appropriately
//
// Parameters:
//
//   - v: The reflect.Value of the struct to marshal
//   - marshalOptions: The protojson.MarshalOptions to use (optional)
//
// Returns:
//
//   - map[string]interface{}: The marshaled struct as a map
//   - error: The error if any
func PrecomputeMarshalByReflection(
	v reflect.Value,
	marshalOptions *protojson.MarshalOptions,
) (map[string]interface{}, error) {
	// Ensure marshalOptions is not nil
	if marshalOptions == nil {
		marshalOptions = &protojson.MarshalOptions{}
	}

	// Dereference pointer if necessary
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Prepare result map
	result := make(map[string]interface{})
	t := v.Type()

	// Handle nested proto.Message fields
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Check if the field can be interfaced
		if field.CanInterface() {
			val := field.Interface()

			// Handle proto.Message fields
			switch fv := val.(type) {
			case proto.Message:
				// Marshal proto.Message to JSON
				data, err := marshalOptions.Marshal(fv)
				if err != nil {
					return nil, err
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
						return nil, err
					}
					result[fieldType.Name] = nested
				} else {
					result[fieldType.Name] = val
				}
			}
		}
	}
	return result, nil
}
