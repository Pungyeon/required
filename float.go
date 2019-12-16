package required

import (
	"encoding/json"
)

// Float32 is a Float32 type, which is required on JSON (un)marshal
type Float32 struct {
	value *float32value
}

type float32value struct {
	value float32
}

// Value will return the inner Float32 type
func (s Float32) Value() float32 {
	return s.value.value
}

// MarshalJSON is an implementation of the json.Marshaler interface
func (s Float32) MarshalJSON() ([]byte, error) {
	if s.value == nil {
		return []byte("null"), nil
	}
	return json.Marshal(s.value)

}

// UnmarshalJSON is an implementation of the json.Unmarhsaler interface
func (s *Float32) UnmarshalJSON(data []byte) error {
	if s.value == nil {
		s.value = &float32value{}
	}
	var v interface{}
	if err := Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		s.value.value = float32(x)
		return nil
	case float32:
		s.value.value = float32(x)
		return nil
	case int32:
		s.value.value = float32(x)
		return nil
	case int64:
		s.value.value = float32(x)
		return nil
	case int:
		s.value.value = float32(x)
		return nil
	default:
		return ErrCannotUnmarshal
	}
}
