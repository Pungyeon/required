package required

import "encoding/json"

// Float32Slice is a required type containing a byte slice value
type Float32Slice struct {
	value []float32
}

// IsValueValid returns whether the contained value has been set
func (s Float32Slice) IsValueValid() error {
	if s.value == nil {
		return ErrEmptyFloatSlice
	}
	return nil
}

// Value will return the inner byte type
func (s Float32Slice) Value() []float32 {
	return s.value
}

// MarshalJSON is an implementation of the json.Marshaler interface
func (s Float32Slice) MarshalJSON() ([]byte, error) {
	if s.Value() == nil {
		return nil, nil
	}
	return json.Marshal(s.value)

}

// UnmarshalJSON is an implementation of the json.Unmarhsaler interface
func (s *Float32Slice) UnmarshalJSON(data []byte) error {
	var v []float32
	if err := Unmarshal(data, &v); err != nil {
		return err
	}
	if len(v) == 0 {
		return ErrEmptyFloatSlice
	}
	s.value = v
	return nil
}

// Float64Slice is a required type containing a byte slice value
type Float64Slice struct {
	value []float64
}

// IsValueValid returns whether the contained value has been set
func (s Float64Slice) IsValueValid() error {
	if s.value == nil {
		return ErrEmptyFloatSlice
	}
	return nil
}

// Value will return the inner byte type
func (s Float64Slice) Value() []float64 {
	return s.value
}

// MarshalJSON is an implementation of the json.Marshaler interface
func (s Float64Slice) MarshalJSON() ([]byte, error) {
	if s.Value() == nil {
		return nil, nil
	}
	return json.Marshal(s.value)

}

// UnmarshalJSON is an implementation of the json.Unmarhsaler interface
func (s *Float64Slice) UnmarshalJSON(data []byte) error {
	var v []float64
	if err := Unmarshal(data, &v); err != nil {
		return err
	}
	if len(v) == 0 {
		return ErrEmptyFloatSlice
	}
	s.value = v
	return nil
}
