package required

import (
	"database/sql"
	"encoding/json"
)

// Float32 is a Float32 type, which is required on JSON (un)marshal
type Float32 struct {
	sql.NullFloat64
}

// IsValueValid returns whether the contained value has been set
func (s Float32) IsValueValid() error {
	if !s.Valid {
		return ErrEmptyFloat
	}
	return nil
}

// Value will return the inner Float32 type
func (s Float32) Value() float32 {
	return float32(s.Float64)
}

// MarshalJSON is an implementation of the json.Marshaler interface
func (s Float32) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return nil, nil
	}
	return json.Marshal(s.Float64)

}

// UnmarshalJSON is an implementation of the json.Unmarhsaler interface
func (s *Float32) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		s.Float64 = float64(float32(x))
		s.Valid = true
		return nil
	default:
		return ErrCannotUnmarshal
	}
}

// Float64 is a Float64 type, which is required on JSON (un)marshal
type Float64 struct {
	sql.NullFloat64
}

// IsValueValid returns whether the contained value has been set
func (s Float64) IsValueValid() error {
	if !s.Valid {
		return ErrEmptyFloat
	}
	return nil
}

// Value will return the inner Float64 type
func (s Float64) Value() float64 {
	return s.Float64
}

// MarshalJSON is an implementation of the json.Marshaler interface
func (s Float64) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return nil, nil
	}
	return json.Marshal(s.Float64)

}

// UnmarshalJSON is an implementation of the json.Unmarhsaler interface
func (s *Float64) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		s.Float64 = float64(x)
		s.Valid = true
		return nil
	default:
		return ErrCannotUnmarshal
	}
}
