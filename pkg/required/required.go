package required

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Nullable struct {
	value interface{}
}

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
	return CheckIfRequired(vo, vo.Type().Name())
}

func CheckIfRequired(vo reflect.Value, field string) error {
	if vo.Kind() != reflect.Struct {
		return nil
	}
	if req, ok := vo.Interface().(Required); ok {
		return validateRequired(req, field)
	}

	return CheckStructFieldsRequired(vo, field)
}

func validateRequired(req Required, field string) error {
	if err := req.IsValueValid(); err != nil {
		return requiredErr{
			err: err,
			msg: field,
		}
	}
	return nil
}

// CheckStructFieldsRequired will inspect the given reflect.Value. If it contains
// a required struct, it will check it's content, if it contains a struct
// it will recursively invoke the function once more, if none of these apply
// nil will be returned.
func CheckStructFieldsRequired(vo reflect.Value, parent string) error {
	for i := 0; i < vo.NumField(); i++ {
		if err := CheckIfRequired(vo.Field(i), childString(parent, vo.Type().Field(i).Name)); err != nil {
			return err
		}
	}
	return nil
}

func childString(parent string, child string) string {
	if parent == "" {
		return child
	}
	return fmt.Sprintf("%s.%s", parent, child)
}
