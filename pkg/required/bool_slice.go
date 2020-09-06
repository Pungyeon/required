package required

import (
	"encoding/json"
)

// BoolSlice is a required type containing a byte slice value
type BoolSlice struct {
	value []bool
}

var _ Required = &BoolSlice{}

func NewBoolSlice(booleans []bool) BoolSlice {
	return BoolSlice{
		value: booleans,
	}
}

// IsValueValid returns whether the contained value has been set
func (s BoolSlice) IsValueValid() error {
	if s.value == nil {
		return ErrEmptyBoolSlice
	}
	return nil
}

// Value will return the inner byte type
func (s BoolSlice) Value() []bool {
	return s.value
}

// MarshalJSON is an implementation of the json.Marshaler interface
func (s BoolSlice) MarshalJSON() ([]byte, error) {
	if s.Value() == nil {
		return nil, nil
	}
	return json.Marshal(s.value)

}

// UnmarshalJSON is an implementation of the json.Unmarshaler interface
func (s *BoolSlice) UnmarshalJSON(data []byte) error {
	var v []bool
	if err := Unmarshal(data, &v); err != nil {
		return err
	}
	if len(v) == 0 {
		return ErrEmptyBoolSlice
	}
	s.value = v
	return nil
}
