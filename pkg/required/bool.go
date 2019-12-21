package required

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Bool is a Bool type, which is required on JSON (un)marshal
type Bool struct {
	value bool
	valid bool
}

var _ Required = Bool{}

// Value will return the inner Bool type
func (s Bool) Value() bool {
	return s.value
}

// IsValueValid returns whether the contained value has been set
func (s Bool) IsValueValid() error {
	if !s.valid {
		return ErrEmpty
	}
	return nil
}

// MarshalJSON is an implementation of the json.Marshaler interface
func (s Bool) MarshalJSON() ([]byte, error) {
	if !s.valid {
		return nil, nil
	}
	return json.Marshal(s.value)

}

// UnmarshalJSON is an implementation of the json.Unmarhsaler interface
func (s *Bool) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := Unmarshal(data, &v); err != nil {
		return err
	}
	fmt.Println(reflect.ValueOf(v))
	switch x := v.(type) {
	case bool:
		fmt.Println("setting the value")
		s.value = bool(x)
		s.valid = true
		return nil
	case map[string]interface{}:
		// This is why embedded fields don't work.
		fmt.Println("wtf do i do?")
		return nil
	default:
		fmt.Println("oh no")
		return ErrCannotUnmarshal
	}
}
