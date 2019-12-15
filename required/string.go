package required

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	// ErrStringEmpty represents an empty required string error
	ErrStringEmpty = errors.New("type of must.String not allowed to be empty")
	// ErrCannotUnmarshal represents an unmarshaling error
	ErrCannotUnmarshal = fmt.Errorf("json: cannot unmarshal into Go value of type must.String")
)

// String is a string type, which is required on JSON (un)marshal
type String struct {
	value string
}

// Value will return the inner string type
func (s *String) Value() string {
	return s.value
}

// MarshalJSON is an implementation of the json.Marshaler interface
func (s String) MarshalJSON() ([]byte, error) {
	if s.Value() == "" {
		return []byte("null"), nil
	}
	return json.Marshal(s.value)

}

// UnmarshalJSON is an implementation of the json.Unmarhsaler interface
func (s *String) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case string:
		if x == "" {
			return ErrStringEmpty
		}
		s.value = x
		return nil
	default:
		return ErrCannotUnmarshal
	}
}
