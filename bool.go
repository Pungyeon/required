package required

import (
	"encoding/json"
)

// Bool is a Bool type, which is required on JSON (un)marshal
type Bool struct {
	value *boolvalue
}

type boolvalue struct {
	value bool
}

// Value will return the inner Bool type
func (s Bool) Value() bool {
	return s.value.value
}

// MarshalJSON is an implementation of the json.Marshaler interface
func (s Bool) MarshalJSON() ([]byte, error) {
	if s.value == nil {
		return []byte("null"), nil
	}
	return json.Marshal(s.value)

}

// UnmarshalJSON is an implementation of the json.Unmarhsaler interface
func (s *Bool) UnmarshalJSON(data []byte) error {
	if s.value == nil {
		s.value = &boolvalue{}
	}
	var v interface{}
	if err := Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case bool:
		s.value.value = bool(x)
		return nil
	default:
		return ErrCannotUnmarshal
	}
}
