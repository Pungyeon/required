package required

import (
	"encoding/json"
	"errors"
)

var (
	// ErrEmpty represents an empty required Int error
	ErrEmpty = errors.New("type of required.Int not allowed to be empty")
)

// Int is a Int type, which is required on JSON (un)marshal
type Int struct {
	value *intvalue
}

type intvalue struct {
	value int
}

// Value will return the inner Int type
func (s Int) Value() int {
	return s.value.value
}

// MarshalJSON is an implementation of the json.Marshaler interface
func (s Int) MarshalJSON() ([]byte, error) {
	if s.value == nil {
		return []byte("null"), nil
	}
	return json.Marshal(s.value)

}

// UnmarshalJSON is an implementation of the json.Unmarhsaler interface
func (s *Int) UnmarshalJSON(data []byte) error {
	if s.value == nil {
		s.value = &intvalue{}
	}
	var v interface{}
	if err := Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		s.value.value = int(x)
		return nil
	case float32:
		s.value.value = int(x)
		return nil
	case int32:
		s.value.value = int(x)
		return nil
	case int64:
		s.value.value = int(x)
		return nil
	case int:
		s.value.value = int(x)
		return nil
	default:
		return ErrCannotUnmarshal
	}
}
