package required

import (
	"database/sql"
	"encoding/json"
)

// Float is a Float type, which is required on JSON (un)marshal
type Float struct {
	sql.NullFloat64
}

// NewFloat returns a valid Float with given value
func NewFloat(value float64) Float {
	return Float{
		sql.NullFloat64{
			Float64: value,
			Valid:   true,
		},
	}
}

// IsValueValid returns whether the contained value has been set
func (s Float) IsValueValid() error {
	if !s.Valid {
		return ErrEmptyFloat
	}
	return nil
}

// Value will return the inner Float type
func (s Float) Value() float64 {
	return s.Float64
}

// MarshalJSON is an implementation of the json.Marshaler interface
func (s Float) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return nil, nil
	}
	return json.Marshal(s.Float64)

}

// UnmarshalJSON is an implementation of the json.Unmarhsaler interface
func (s *Float) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		s.Float64 = x
		s.Valid = true
		return nil
	default:
		return ErrCannotUnmarshal
	}
}
