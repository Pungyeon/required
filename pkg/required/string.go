package required

import (
	"encoding/json"
)

// String is a string type, which is required on JSON (un)marshal
type String struct {
	value string
}

// IsValueValid returns whether the contained value has been set
func (s String) IsValueValid() error {
	if s.value == "" {
		return ErrEmpty
	}
	return nil
}

// Value will return the inner string type
func (s String) Value() string {
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
			return ErrEmpty
		}
		s.value = x
		return nil
	default:
		return ErrCannotUnmarshal
	}
}
