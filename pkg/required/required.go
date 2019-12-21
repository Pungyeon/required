package required

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
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

// Required is an interface which will enable the require.Unmarshal parser,
// to check whether a given object / interface has a valid contained value.
type Required interface {
	IsValueValid() error
}

// ReturnIfError will iterate over a variadac error and return
// an error if the given value is not nil
func ReturnIfError(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

// Unmarshal is a wrapping function of the json.Unmarshal function
func Unmarshal(data []byte, v interface{}) error {
	return ReturnIfError(
		json.Unmarshal(data, v),
		CheckValues(v),
	)
}

// CheckValues will check the values of a given interface and ensure
// that if it contains a required struct, that the required values
// are not empty
func CheckValues(v interface{}) error {
	vo := reflect.ValueOf(v)
	for vo.Kind() == reflect.Ptr {
		vo = vo.Elem()
	}
	return CheckStructIsRequired(vo)
}

// CheckStructIsRequired will inspect the given reflect.Value. If it contains
// a required struct, it will check it's content, if it contains a struct
// it will recursively invoke the function once more, if none of these apply
// nil will be returned.
func CheckStructIsRequired(vo reflect.Value) error {
	if vo.Kind() != reflect.Struct {
		return nil
	}
	for i := 0; i < vo.NumField(); i++ {
		vtf := vo.Field(i)
		if req, ok := vtf.Interface().(Required); ok {
			return req.IsValueValid()
		}
		if vtf.Kind() == reflect.Struct {
			if err := CheckStructIsRequired(vtf); err != nil {
				return err
			}
		}
	}
	return nil
}
