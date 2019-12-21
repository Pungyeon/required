package required

import (
	"encoding/json"
)

// Int is a Int type, which is required on JSON (un)marshal
type Int struct {
	value int
	valid bool
}

// IsValueValid returns whether the contained value has been set
func (s Int) IsValueValid() error {
	if !s.valid {
		return ErrEmpty
	}
	return nil
}

// Value will return the inner Int type
func (s Int) Value() int {
	return s.value
}

// MarshalJSON is an implementation of the json.Marshaler interface
func (s Int) MarshalJSON() ([]byte, error) {
	if !s.valid {
		return nil, nil
	}
	return json.Marshal(s.value)

}

// UnmarshalJSON is an implementation of the json.Unmarhsaler interface
func (s *Int) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		s.value = int(x)
		s.valid = true
		return nil
	default:
		return ErrCannotUnmarshal
	}
}
