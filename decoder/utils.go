package decoder

import (
	"bytes"
	"io"
	"strings"
)

// ToReader converts an any to an io.Reader
//
// Parameters:
//
//   - reader: The any to convert
//
// Returns:
//
// - io.Reader: The converted io.Reader
// - error: Error if the conversion fails
func ToReader(reader any) (io.Reader, error) {
	switch v := reader.(type) {
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
