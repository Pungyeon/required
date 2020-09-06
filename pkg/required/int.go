package required

import (
	"encoding/json"
)

// Int is a Int type, which is required on JSON (un)marshal
type Int struct {
	Nullable
}

var _ Required = Int{}

// NewInt returns a valid Int with given value
func NewInt(value int) Int {
	return Int{
		Nullable{
			value: value,
		},
	}
}

// IsValueValid returns whether the contained value has been set
func (s Int) IsValueValid() error {
	if s.value == nil {
		return ErrEmptyInt
	}
	return nil
}

// Value will return the inner Int type
func (s Int) Value() int {
	return s.value.(int)
}

// MarshalJSON is an implementation of the json.Marshaler interface
func (s Int) MarshalJSON() ([]byte, error) {
	if err := s.IsValueValid(); err != nil {
		return nil, err
	}
	return json.Marshal(s.Value())

}

// UnmarshalJSON is an implementation of the json.Unmarhsaler interface
func (s *Int) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		s.value = int(x)
		return nil
	default:
		return ErrCannotUnmarshal
	}
}
