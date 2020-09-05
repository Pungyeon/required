package required

import (
	"errors"
	"fmt"
)

var (
	ErrCannotUnmarshal  = fmt.Errorf("json: cannot unmarshal given value")
	ErrEmpty            = errors.New("type of required.Required not allowed to be empty")
	ErrEmptyBool        = errors.New("type of required.Bool not allowed to be empty")
	ErrEmptyBoolSlice   = errors.New("type of required.BoolSlice not allowed to be empty")
	ErrEmptyByteSlice   = errors.New("type of required.ByteSlice not allowed to be empty")
	ErrEmptyFloatSlice  = errors.New("type of required.FloatSlice not allowed to be empty")
	ErrEmptyFloat       = errors.New("type of required.Float not allowed to be empty")
	ErrEmptyIntSlice    = errors.New("type of required.IntSlice not allowed to be empty")
	ErrEmptyInt         = errors.New("type of required.Int not allowed to be empty")
	ErrEmptyStringSlice = errors.New("type of required.StringSlice not allowed to be empty")
	ErrEmptyString      = errors.New("type of required.String not allowed to be empty")
)

type requiredErr struct {
	err error
	msg string
}

func (err requiredErr) Error() string {
	return fmt.Sprintf("%s: %v", err.msg, err.err)
}

// IsRequiredErr will type check the given error as a requiredErr
// returning a boolean on whether the type check was successful
func IsRequiredErr(err error) bool {
	_, ok := err.(requiredErr)
	return ok
}
