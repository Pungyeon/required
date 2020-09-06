package required

import (
	"encoding/json"
)

// Float is a Float type, which is required on JSON (un)marshal
type Float struct {
	Nullable
}

// NewFloat returns a valid Float with given value
func NewFloat(value float64) Float {
	return Float{
		Nullable{
			value: value,
		},
	}
}

// IsValueValid returns whether the contained value has been set
func (s Float) IsValueValid() error {
	if s.value == nil {
		return ErrEmptyFloat
	}
	return nil
}

// Value will return the inner Float type
func (s Float) Value() float64 {
	return s.value.(float64)
}

// MarshalJSON is an implementation of the json.Marshaler interface
func (s Float) MarshalJSON() ([]byte, error) {
	if err := s.IsValueValid(); err != nil {
		return nil, err
	}
	return json.Marshal(s.Value())

}

// UnmarshalJSON is an implementation of the json.Unmarhsaler interface
func (s *Float) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		s.value = x
		return nil
	default:
		return ErrCannotUnmarshal
	}
}
