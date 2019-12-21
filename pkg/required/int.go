package required

import (
	"database/sql"
	"encoding/json"
)

// Int is a Int type, which is required on JSON (un)marshal
type Int struct {
	sql.NullInt64
}

// IsValueValid returns whether the contained value has been set
func (s Int) IsValueValid() error {
	if !s.Valid {
		return ErrEmpty
	}
	return nil
}

// Value will return the inner Int type
func (s Int) Value() int {
	return int(s.Int64)
}

// MarshalJSON is an implementation of the json.Marshaler interface
func (s Int) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return nil, nil
	}
	return json.Marshal(s.Int64)

}

// UnmarshalJSON is an implementation of the json.Unmarhsaler interface
func (s *Int) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		s.Int64 = int64(x)
		s.Valid = true
		return nil
	default:
		return ErrCannotUnmarshal
	}
}
