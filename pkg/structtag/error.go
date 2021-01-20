package structtag

import (
	"errors"
	"fmt"
)

var (
	errRequiredField = errors.New("RequiredInterface field missing")
)

func IsRequiredErr(err error) bool {
	_, ok := err.(requiredErr)
	return ok
}

type requiredErr struct {
	err   error
	field string
}

func (err requiredErr) Error() string {
	return fmt.Sprintf("%v: %s", err.err, err.field)
}

var _ error = &requiredErr{}
