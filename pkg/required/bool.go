package required

import (
	"database/sql"
	"encoding/json"
)

// Bool is a Bool type, which is required on JSON (un)marshal
type Bool struct {
	sql.NullBool
}

var _ Required = Bool{}

// NewBool returns a valid Bool with given value
func NewBool(value bool) Bool {
	return Bool{
		NullBool: sql.NullBool{
			Bool: value,
			Valid: true,
		},
	}
}

// Value will return the inner Bool type
func (s Bool) Value() bool {
	return s.Bool
}

// IsValueValid returns whether the contained value has been set
func (s Bool) IsValueValid() error {
	if !s.Valid {
		return ErrEmptyBool
	}
	return nil
}

// MarshalJSON is an implementation of the json.Marshaler interface
func (s Bool) MarshalJSON() ([]byte, error) {
	if err := s.IsValueValid(); err != nil {
		return nil, err
	}
	return json.Marshal(s.Bool)

}

// UnmarshalJSON is an implementation of the json.Unmarhsaler interface
func (s *Bool) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case bool:
		s.Bool = bool(x)
		s.Valid = true
		return nil
	default:
		return ErrCannotUnmarshal
	}
}
