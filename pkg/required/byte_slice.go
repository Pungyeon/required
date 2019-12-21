package required

import "encoding/json"

// ByteSlice is a required type containing a byte slice value
type ByteSlice struct {
	value []byte
}

// Value will return the inner byte type
func (s ByteSlice) Value() []byte {
	return s.value
}

// IsValueValid returns whether the contained value has been set
func (s ByteSlice) IsValueValid() error {
	if s.value == nil {
		return ErrEmpty
	}
	return nil
}

// MarshalJSON is an implementation of the json.Marshaler interface
func (s ByteSlice) MarshalJSON() ([]byte, error) {
	if s.Value() == nil {
		return nil, nil
	}
	return json.Marshal(s.value)

}

// UnmarshalJSON is an implementation of the json.Unmarhsaler interface
func (s *ByteSlice) UnmarshalJSON(data []byte) error {
	var v []byte
	if err := Unmarshal(data, &v); err != nil {
		return err
	}
	if len(v) == 0 {
		return ErrEmpty
	}
	s.value = v
	return nil
}
