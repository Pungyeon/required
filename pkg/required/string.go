package required

import (
	"encoding/json"
)

// String is a string type, which is required on JSON (un)marshal
type String struct {
	Nullable
}

// NewString returns a valid String with given value
func NewString(str string) String {
	return String{
		Nullable{
			value: str,
		},
	}
}

// IsValueValid returns whether the contained value has been set
func (s String) IsValueValid() error {
	if s.value == nil {
		return ErrEmptyString
	}
	if s.Value() == "" {
		return ErrEmptyString
	}
	return nil
}

// Value will return the inner string type
func (s String) Value() string {
	return s.value.(string)
}

// MarshalJSON is an implementation of the json.Marshaler interface
func (s String) MarshalJSON() ([]byte, error) {
	if err := s.IsValueValid(); err != nil {
		return nil, err
	}
	return json.Marshal(s.Value())

}

// UnmarshalJSON is an implementation of the json.Unmarhsaler interface
func (s *String) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case string:
		s.value = x
		return nil
	default:
		return ErrCannotUnmarshal
	}
}
