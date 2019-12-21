package required

import "encoding/json"

// BoolSlice is a required type containing a byte slice value
type BoolSlice struct {
	value []bool
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

// UnmarshalJSON is an implementation of the json.Unmarhsaler interface
func (s *BoolSlice) UnmarshalJSON(data []byte) error {
	var v []bool
	if err := Unmarshal(data, &v); err != nil {
		return err
	}
	if len(v) == 0 {
		return ErrEmpty
	}
	s.value = v
	return nil
}
