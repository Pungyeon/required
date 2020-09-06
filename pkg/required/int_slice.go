package required

import "encoding/json"

// IntSlice is a required type containing a int slice value
type IntSlice struct {
	value []int
}

// NewIntSlice returns a valid IntSlice with given value
func NewIntSlice(ints []int) IntSlice {
	return IntSlice{
		value: ints,
	}
}

// IsValueValid returns whether the contained value has been set
func (s IntSlice) IsValueValid() error {
	if s.value == nil {
		return ErrEmptyIntSlice
	}
	return nil
}

// Value will return the inner int type
func (s IntSlice) Value() []int {
	return s.value
}

// MarshalJSON is an implementation of the json.Marshaler interface
func (s IntSlice) MarshalJSON() ([]byte, error) {
	if err := s.IsValueValid(); err != nil {
		return nil, err
	}
	return json.Marshal(s.value)
}

// UnmarshalJSON is an implementation of the json.Unmarhsaler interface
func (s *IntSlice) UnmarshalJSON(data []byte) error {
	var v []int
	if err := Unmarshal(data, &v); err != nil {
		return err
	}
	if len(v) == 0 {
		return ErrEmptyIntSlice
	}
	s.value = v
	return nil
}
