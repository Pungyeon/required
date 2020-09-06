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
		checkValues(v),
	)
}

// checkValues will check the values of a given interface and ensure
// that if it contains a required struct, that the required values
// are not empty
func checkValues(v interface{}) error {
	vo := reflect.ValueOf(v)
	return checkIfRequired(vo, vo.Type().Name())
}

// checkIfRequired will ensure that the given value is valid
// if it is a structure which fulfils the Required interface
func checkIfRequired(vo reflect.Value, field string) error {
	for vo.Kind() == reflect.Ptr {
		vo = vo.Elem()
	}
	if vo.Kind() != reflect.Struct {
		return nil
	}
	if req, ok := vo.Interface().(Required); ok {
		return validateRequired(req, field)
	}
	return checkStructFieldsRequired(vo, field)
}

// validateRequired is a wrapper function for invoking the
// Required interface and returning a detailed error, if
// the value is invalid
func validateRequired(req Required, field string) error {
	if err := req.IsValueValid(); err != nil {
		return requiredErr{
			err: err,
			msg: field,
		}
	}
	return nil
}

// checkStructFieldsRequired will inspect the given reflect.Value. If it contains
// a required struct, it will check it's content, if it contains a struct
// it will recursively invoke the function once more, if none of these apply
// nil will be returned.
func checkStructFieldsRequired(vo reflect.Value, parent string) error {
	for i := 0; i < vo.NumField(); i++ {
		if err := checkIfRequired(
			vo.Field(i),
			childString(parent, vo.Type().Field(i).Name),
		); err != nil {
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
