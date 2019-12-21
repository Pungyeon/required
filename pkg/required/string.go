package required

import (
	"database/sql"
	"encoding/json"
)

// String is a string type, which is required on JSON (un)marshal
type String struct {
	sql.NullString
}

// IsValueValid returns whether the contained value has been set
func (s String) IsValueValid() error {
	if s.String == "" {
		return ErrEmpty
	}
	return nil
}

// Value will return the inner string type
func (s String) Value() string {
	return s.String
}

// MarshalJSON is an implementation of the json.Marshaler interface
func (s String) MarshalJSON() ([]byte, error) {
	if s.Value() == "" {
		return []byte("null"), nil
	}
	return json.Marshal(s.String)

}

// UnmarshalJSON is an implementation of the json.Unmarhsaler interface
func (s *String) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case string:
		s.String = x
		return nil
	default:
		return ErrCannotUnmarshal
	}
}
