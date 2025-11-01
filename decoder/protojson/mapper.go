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
	// Mapper is the struct to hold precomputed marshal by reflection functions
	Mapper struct {
		reflectType        reflect.Type
		isProtoMessage     bool
		regularFields      map[string]struct{}
		protoMessageFields map[string]struct{}
		jsonFieldNames     map[string]string
		nestedStructs      map[string]*Mapper
	}
)

// NewMapper creates a new Mapper instance
//
// Parameters:
//
// - destinationInstance: the destination instance to create the mapper for
//
// Returns:
//
// - *Mapper: the new Mapper instance
// - error: the error if any
func NewMapper(
	destinationInstance any,
) (*Mapper, error) {
	// Check if the destination instance is nil
	if destinationInstance == nil {
		return nil, ErrNilDestinationInstance
	}

	// Check if the destination is a proto.Message
	_, ok := destinationInstance.(proto.Message)
	if ok {
		return &Mapper{
			isProtoMessage: true,
		}, nil
	}

	// Create the maps to hold field information
	regularFields := make(map[string]struct{})
	protoMessageFields := make(map[string]struct{})
	nestedStructs := make(map[string]*Mapper)
	jsonFieldNames := make(map[string]string)

	// Get the reflect type and value of the destination instance
	reflectType := goreflect.GetDereferencedType(destinationInstance)
	reflectValue := goreflect.GetDereferencedValue(destinationInstance)

	// Check if the type is a struct
	for i := 0; i < reflectValue.NumField(); i++ {
		// Get the field and its type
		structField := reflectType.Field(i)
		fieldValue := reflectValue.Field(i)
		fieldName := structField.Name

		// Check if the field can be set
		if !fieldValue.CanSet() {
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

		// Store the JSON field name
		jsonFieldNames[fieldName] = jsonFieldName

		// Get the field interface
		fieldValueInterface := fieldValue.Interface()

		// Check if the field is a proto.Message
		if _, protoOk := fieldValueInterface.(proto.Message); protoOk {
			// Store the proto.Message field
			protoMessageFields[fieldName] = struct{}{}
			continue
		}

		// Dereference pointer if necessary
		if fieldValue.Kind() == reflect.Ptr {
			fieldValue = fieldValue.Elem()
		}

		// Check if the field is a struct
		if fieldValue.Kind() == reflect.Struct {
			// Create a nested mapper for the struct field
			nestedMapper, nestedErr := NewMapper(fieldValueInterface)
			if nestedErr != nil {
				return nil, nestedErr
			}
			nestedStructs[fieldName] = nestedMapper
			continue
		}

		// Regular field
		regularFields[fieldName] = struct{}{}
	}
	return &Mapper{
		reflectType:        reflectType,
		isProtoMessage:     false,
		regularFields:      regularFields,
		protoMessageFields: protoMessageFields,
		nestedStructs:      nestedStructs,
	}, nil
}

// UnmarshalByReflection unmarshal JSON data into a destination using reflection
//
// Parameters:
//
//   - body: The JSON data to unmarshal
//   - dest: The destination to unmarshal the JSON data into
//   - unmarshalOptions: Options for unmarshalling proto messages (optional, can be nil)
//
// Returns:
//
//   - error: The error if any
func (m *Mapper) UnmarshalByReflection(
	body []byte,
	dest any,
	unmarshalOptions *protojson.UnmarshalOptions,
) error {
	// Check if the mapper is nil
	if m == nil {
		return ErrNilMapper
	}

	// Check if the destination is nil
	if dest == nil {
		return ErrNilDestination
	}

	// Check if the unmarshal options are nil, if so initialize them
	if unmarshalOptions == nil {
		unmarshalOptions = &protojson.UnmarshalOptions{}
	}

	// Check if the destination is a proto.Message
	if m.isProtoMessage {
		// Ensure the destination is a proto.Message
		parsedProtoMessage, ok := dest.(proto.Message)
		if !ok {
			return ErrDestinationNotProtoMessage
		}

		// Unmarshal directly into the proto.Message
		return unmarshalOptions.Unmarshal(
			body,
			parsedProtoMessage,
		)
	}

	// Initialize the map to hold intermediate JSON data
	var tempDest map[string]any
	if unmarshalErr := json.Unmarshal(body, &tempDest); unmarshalErr != nil {
		return unmarshalErr
	}

	// Get the reflect value of the destination
	reflectValue := goreflect.GetDereferencedValue(dest)

	// Decode nested proto messages
	mappedBody := make(map[string]any)
	for i := 0; i < reflectValue.NumField(); i++ {
		// Get the field and its type
		structField := m.reflectType.Field(i)
		fieldName := structField.Name
		fieldType := structField.Type
		fieldValue := reflectValue.Field(i)

		// Get the JSON tag of the field
		jsonFieldName, jsonFielNameOk := m.jsonFieldNames[fieldName]
		if !jsonFielNameOk {
			// If no JSON field name is found, skip the field, it means it cannot be set
			continue
		}

		// Get the corresponding body field
		bodyField, ok := tempDest[jsonFieldName]
		if !ok {
			continue
		}

		// Check if the field is a regular field
		if _, regularOk := m.regularFields[fieldName]; regularOk {
			// Directly map the body field
			mappedBody[fieldName] = bodyField
			continue
		}

		// Check if the field is a proto.Message
		if _, protoOk := m.protoMessageFields[fieldName]; protoOk {
			// Marshal the body field to JSON
			marshaledField, err := json.Marshal(bodyField)
			if err != nil {
				return err
			}

			// Create a new instance of the proto.Message if it's a pointer and is nil
			if fieldValue.Kind() == reflect.Ptr {
				fieldValue.Set(reflect.New(fieldType.Elem()))
			}

			// Get the interface of the field value
			fieldValueInterface := fieldValue.Interface()

			// Parse field value interface as proto.Message
			protoMessage, parsedProtoOk := fieldValueInterface.(proto.Message)
			if !parsedProtoOk {
				return fmt.Errorf(ErrFieldNotProtoMessage, fieldName)
			}

			// Unmarshal the JSON into the proto.Message field
			if unmarshalErr := unmarshalOptions.Unmarshal(
				marshaledField,
				protoMessage,
			); unmarshalErr != nil {
				return unmarshalErr
			}

			// Map the proto.Message field
			mappedBody[fieldName] = fieldValueInterface
			continue
		}

		// Check if the field is a nested struct
		if nestedMapper, nestedOk := m.nestedStructs[fieldName]; nestedOk {
			// Marshal the body field to JSON
			marshaledNestedBody, err := json.Marshal(bodyField)
			if err != nil {
				return err
			}

			// Get the interface of the field value
			fieldValueInterface := fieldValue.Addr().Interface()

			// Unmarshal the body field by reflection
			if nestedErr := nestedMapper.UnmarshalByReflection(marshaledNestedBody, fieldValueInterface, unmarshalOptions); nestedErr != nil {
				return nestedErr
			}

			// Set the mapped nested struct
			mappedBody[fieldName] = fieldValueInterface
			continue
		}

		// The field type is not handled, return an error
		return fmt.Errorf(ErrFieldNotHandled, fieldName)
	}

	// Reflect the destination value
	destValue := goreflect.GetDereferencedValue(dest)

	// Map the body to the struct field
	return goreflect.ValueToReflectStruct(mappedBody, destValue, false)
}
