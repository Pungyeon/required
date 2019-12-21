package required

import (
	"encoding/json"
)

// StringSlice is a required type containing a string slice value
type StringSlice struct {
	value []string
}

// Value will return the inner string type
func (s StringSlice) Value() []string {
	return s.value
}

// IsValueValid returns whether the contained value has been set
func (s StringSlice) IsValueValid() error {
	if s.value == nil {
		return ErrEmpty
	}
	return nil
}

// MarshalJSON is an implementation of the json.Marshaler interface
func (s StringSlice) MarshalJSON() ([]byte, error) {
	if s.Value() == nil {
		return nil, nil
	}
	return json.Marshal(s.value)

}

// UnmarshalJSON is an implementation of the json.Unmarhsaler interface
func (s *StringSlice) UnmarshalJSON(data []byte) error {
	var v []string
	if err := Unmarshal(data, &v); err != nil {
		return err
	}
	if len(v) == 0 {
		return ErrEmpty
	}
	s.value = v
	return nil
}
