package required

import (
	"encoding/json"
)

// Float32 is a Float32 type, which is required on JSON (un)marshal
type Float32 struct {
	value float32
	valid bool
}

// IsValueValid returns whether the contained value has been set
func (s Float32) IsValueValid() error {
	if !s.valid {
		return ErrEmpty
	}
	return nil
}

// Value will return the inner Float32 type
func (s Float32) Value() float32 {
	return s.value
}

// MarshalJSON is an implementation of the json.Marshaler interface
func (s Float32) MarshalJSON() ([]byte, error) {
	if !s.valid {
		return nil, nil
	}
	return json.Marshal(s.value)

}

// UnmarshalJSON is an implementation of the json.Unmarhsaler interface
func (s *Float32) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		s.value = float32(x)
		s.valid = true
		return nil
	default:
		return ErrCannotUnmarshal
	}
}

// Float64 is a Float64 type, which is required on JSON (un)marshal
type Float64 struct {
	value float64
	valid bool
}

// IsValueValid returns whether the contained value has been set
func (s Float64) IsValueValid() error {
	if !s.valid {
		return ErrEmpty
	}
	return nil
}

// Value will return the inner Float64 type
func (s Float64) Value() float64 {
	return s.value
}

// MarshalJSON is an implementation of the json.Marshaler interface
func (s Float64) MarshalJSON() ([]byte, error) {
	if !s.valid {
		return nil, nil
	}
	return json.Marshal(s.value)

}

// UnmarshalJSON is an implementation of the json.Unmarhsaler interface
func (s *Float64) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		s.value = float64(x)
		s.valid = true
		return nil
	default:
		return ErrCannotUnmarshal
	}
}
