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
		protoMessageFields        map[string]struct{}
		regularFields             map[string]struct{}
		jsonFieldNames           map[string]string
		nestedStructs map[string]*Mapper
	}
)

// NewMapper creates a new protoJSON mapper
// 
// Parameters:
// 
//  - structInstance: instance of the struct to create the mapper from
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
	reflection := goreflect.NewReflection(structInstance)
	reflectedType := reflection.GetReflectedType()
	reflectedValue := reflection.GetReflectedValue()
	
	// Prepare the different maps
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
		jsonTag, err := gostringsjson.GetJSONTag(structField, fieldName)
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
			nestedMapper, err := NewMapper(fieldInterfaceValue)
			if err != nil {
				return nil, err
			}
			nestedStructs[fieldName] = nestedMapper
		}
	}
	return &Mapper{
		protoMessageFields: protoMessageFields,
		regularFields:      regularFields,
		jsonFieldNames:    jsonFieldNames,
		nestedStructs: nestedStructs,
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
	//	return e.jsonEncoder.Encode(precomputedMarshal)/ Ensure marshalOptions is not nil
	if marshalOptions == nil {
		marshalOptions = &protojson.MarshalOptions{}
	}

	// Reflect on the instance to get its fields
	reflectValue := goreflect.GetDereferencedValue(body)
	reflectType := goreflect.GetDereferencedType(body)
	
	// Prepare result map
	result := make(map[string]any)
	
	// Handle nested proto.Message fields
	for i := 0; i < reflectValue.NumField(); i++ {
		field := reflectValue.Field(i)
		fieldType := reflectType.Field(i)
		fieldName := fieldType.Name
		
		// Try to get the JSON field name
		jsonFieldName, ok := m.jsonFieldNames[fieldName]
		if !ok {
			// Continue, it means the field is not mapped (unexported or ignored)
			continue
		}
		
		// Get the field value
		fieldValue := field.Interface()
		
		// Check if the field is a regular field 
		if _, ok := m.regularFields[fieldName]; ok {
			// Process the field
			result[jsonFieldName] = fieldValue
			continue
		}
		
		// Check if the field is a nested struct
		if nestedMapper, ok := m.nestedStructs[fieldName]; ok {
			// Recursively process the nested struct
			nestedResult, err := nestedMapper.PrecomputeMarshalByReflection(
				fieldValue,
				marshalOptions,
			)
			if err != nil {
				return nil, err
			}
			result[jsonFieldName] = nestedResult
			continue
		}
			
		// Check if the field is a proto.Message field
		if _, ok := m.protoMessageFields[fieldName]; ok {
			// Get the field value as proto.Message
			fieldValue, ok := fieldValue.(proto.Message)
			if !ok {
				return nil, fmt.Errorf(ErrFieldNotProtoMessage, fieldName)
			}
			
			// Marshal proto.Message to JSON
			data, err := marshalOptions.Marshal(fieldValue)
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
