package protojson

import (
	"encoding/json"
	"fmt"
	"reflect"

	goreflect "github.com/ralvarezdev/go-reflect"
	gostringsjson "github.com/ralvarezdev/go-strings/json"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type (
	// Mapper is the protoJSON mapper struct
	Mapper struct {
		reflectType        reflect.Type
		optionalFields     map[string]struct{}
		protoMessageFields map[string]struct{}
		regularFields      map[string]struct{}
		jsonFieldNames     map[string]string
		nestedStructs      map[string]*Mapper
	}
)

// NewMapper creates a new protoJSON mapper
//
// Parameters:
//
//   - structInstance: instance of the struct to create the mapper from
//
// Returns:
//
// - *Mapper: instance of the mapper
// - error: error if the struct instance is nil
func NewMapper(structInstance any) (*Mapper, error) {
	// Check if the struct instance is nil
	if structInstance == nil {
		return nil, ErrNilStructInstance
	}

	// Reflection of data
	reflectedType := goreflect.GetDereferencedType(structInstance)
	reflectedValue := goreflect.GetDereferencedValue(structInstance)

	// Prepare the different maps
	optionalFields := make(map[string]struct{})
	protoMessageFields := make(map[string]struct{})
	regularFields := make(map[string]struct{})
	jsonFieldNames := make(map[string]string)
	nestedStructs := make(map[string]*Mapper)

	// Handle nested proto.Message fields
	for i := 0; i < reflectedType.NumField(); i++ {
		// Get the field type through reflection
		structField := reflectedType.Field(i)
		fieldValue := reflectedValue.Field(i)
		fieldType := structField.Type
		fieldName := structField.Name

		// Check if the field can be interfaced
		if !fieldValue.CanInterface() {
			continue
		}

		// Get the JSON tag of the field
		jsonTag, err := gostringsjson.GetJSONTag(&structField, fieldName)
		if err != nil {
			return nil, err
		}

		// Get the JSON field name from the tag
		jsonFieldName, err := gostringsjson.GetJSONTagName(jsonTag, fieldName)
		if err != nil {
			return nil, err
		}
		
		// Check if the field is optional
		if gostringsjson.IsJSONFieldOptional(jsonTag){
			optionalFields[fieldName] = struct{}{}
		}

		// Store the JSON field name
		jsonFieldNames[fieldName] = jsonFieldName

		// Get the field interface value
		fieldInterfaceValue := fieldValue.Interface()

		// Handle proto.Message fields
		switch fieldInterfaceValue.(type) {
		case proto.Message:
			// Set the field as a protoMessageField
			protoMessageFields[fieldName] = struct{}{}
		default:
			if fieldType.Kind() != reflect.Struct {
				// Store as regular field
				regularFields[fieldName] = struct{}{}
				break
			}

			// Recursively handle nested structs
			nestedMapper, mapperErr := NewMapper(fieldInterfaceValue)
			if mapperErr != nil {
				return nil, mapperErr
			}
			nestedStructs[fieldName] = nestedMapper
		}
	}
	return &Mapper{
		reflectType:        reflectedType,
		optionalFields:     optionalFields,
		protoMessageFields: protoMessageFields,
		regularFields:      regularFields,
		jsonFieldNames:     jsonFieldNames,
		nestedStructs:      nestedStructs,
	}, nil
}

// PrecomputeMarshalByReflection marshals a struct to a map[string]any using reflection, handling nested
// proto.Message fields appropriately
//
// Parameters:
//
//   - body: The struct to marshal
//   - marshalOptions: The protojson.MarshalOptions to use (optional)
//
// Returns:
//
//   - map[string]any: The marshaled struct as a map
//   - error: The error if any
func (m *Mapper) PrecomputeMarshalByReflection(
	body any,
	marshalOptions *protojson.MarshalOptions,
) (map[string]any, error) {
	// Check if the mapper is nil
	if m == nil {
		return nil, ErrNilMapper
	}

	// Check if body is nil
	if body == nil {
		return nil, ErrNilBody
	}

	// Check if the marshal options are nil, if so create a default one
	if marshalOptions == nil {
		marshalOptions = &protojson.MarshalOptions{}
	}

	// Reflect on the instance to get its fields
	reflectValue := goreflect.GetDereferencedValue(body)

	// Prepare result map
	result := make(map[string]any)

	// Handle nested proto.Message fields
	for i := 0; i < reflectValue.NumField(); i++ {
		// Get the field type through reflection
		structField := m.reflectType.Field(i)
		fieldValue := reflectValue.Field(i)
		fieldName := structField.Name

		// Try to get the JSON field name
		jsonFieldName, ok := m.jsonFieldNames[fieldName]
		if !ok {
			// If not found, skip the field, it means it cannot be interfaced
			continue
		}

		// Get the field value
		fieldValueInterface := fieldValue.Interface()
		
		// Check if the field is optional and zero value
		if _, optionalOk := m.optionalFields[fieldName]; optionalOk && fieldValue.IsZero() {
			continue
		}

		// Check if the field is a regular field
		if _, regularOk := m.regularFields[fieldName]; regularOk {
			// Process the field
			result[jsonFieldName] = fieldValueInterface
			continue
		}

		// Check if the field is a nested struct
		if nestedMapper, nestedOk := m.nestedStructs[fieldName]; nestedOk {
			// Recursively process the nested struct
			nestedResult, err := nestedMapper.PrecomputeMarshalByReflection(
				fieldValueInterface,
				marshalOptions,
			)
			if err != nil {
				return nil, err
			}
			result[jsonFieldName] = nestedResult
			continue
		}

		// Check if the field is a proto.Message field
		if _, protoMessageOk := m.protoMessageFields[fieldName]; protoMessageOk {
			// Get the field value as proto.Message
			protoMessage, protoOk := fieldValueInterface.(proto.Message)
			if !protoOk {
				return nil, fmt.Errorf(ErrFieldNotProtoMessage, fieldName)
			}

			// Marshal proto.Message to JSON
			data, err := marshalOptions.Marshal(protoMessage)
			if err != nil {
				return nil, err
			}

			// Store as json.RawMessage to avoid double encoding
			result[jsonFieldName] = json.RawMessage(data)
			continue
		}

		// The field type is not handled, return an error
		return nil, fmt.Errorf(ErrFieldNotHandled, fieldName)
	}
	return result, nil
}
