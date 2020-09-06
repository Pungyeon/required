package required

import "encoding/json"

// FloatSlice is a required type containing a byte slice value
type FloatSlice struct {
	value []float64
}

var _ Required = FloatSlice{}

// NewFloatSlice returns a valid FloatSlice with given value
func NewFloatSlice(floats []float64) FloatSlice {
	return FloatSlice{
		value: floats,
	}
}

// IsValueValid returns whether the contained value has been set
func (s FloatSlice) IsValueValid() error {
	if s.value == nil {
		return ErrEmptyFloatSlice
	}
	return nil
}

// Value will return the inner byte type
func (s FloatSlice) Value() []float64 {
	return s.value
}

// MarshalJSON is an implementation of the json.Marshaler interface
func (s FloatSlice) MarshalJSON() ([]byte, error) {
	if err := s.IsValueValid(); err != nil {
		return nil, err
	}
	return json.Marshal(s.value)

}

// UnmarshalJSON is an implementation of the json.Unmarhsaler interface
func (s *FloatSlice) UnmarshalJSON(data []byte) error {
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
