package decoder

import (
	"bytes"
	"io"
	"strings"
)

// ToReader converts an interface{} to an io.Reader
//
// Parameters:
//
//   - any: The interface{} to convert (can be io.Reader, string, or []byte)
//
// Returns:
//
// - io.Reader: The converted io.Reader
// - error: Error if the conversion fails
func ToReader(any interface{}) (io.Reader, error) {
	switch v := any.(type) {
	case io.Reader:
		return v, nil
	case string:
		return strings.NewReader(v), nil
	case []byte:
		return bytes.NewReader(v), nil
	default:
		return nil, ErrInvalidInstance
	}
}
