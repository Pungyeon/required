package required

import "encoding/json"

// IntSlice is a required type containing a int slice value
type IntSlice struct {
	value []int
}

// Value will return the inner int type
func (s IntSlice) Value() []int {
	return s.value
}

// MarshalJSON is an implementation of the json.Marshaler interface
func (s IntSlice) MarshalJSON() ([]byte, error) {
	if s.Value() == nil {
		return nil, nil
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
		return ErrEmpty
	}
	s.value = v
	return nil
}
