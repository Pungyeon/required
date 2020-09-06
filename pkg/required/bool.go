package required

import (
	"encoding/json"
)

// Bool is a Bool type, which is required on JSON (un)marshal
type Bool struct {
	Nullable
}

var _ Required = Bool{}

// NewBool returns a valid Bool with given value
func NewBool(value bool) Bool {
	return Bool{
		Nullable{
			value: value,
		},
	}
}

// Value will return the inner Bool type
func (s Bool) Value() bool {
	return s.value.(bool)
}

// IsValueValid returns whether the contained value has been set
func (s Bool) IsValueValid() error {
	if s.value == nil {
		return ErrEmptyBool
	}
	return nil
}

// MarshalJSON is an implementation of the json.Marshaler interface
func (s Bool) MarshalJSON() ([]byte, error) {
	if err := s.IsValueValid(); err != nil {
		return nil, err
	}
	return json.Marshal(s.Value())

}

// UnmarshalJSON is an implementation of the json.Unmarhsaler interface
func (s *Bool) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case bool:
		s.value = x
		return nil
	default:
		return ErrCannotUnmarshal
	}
}
