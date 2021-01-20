package required

import (
	"encoding/json"
	"reflect"
)

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

// Unmarshal is a wrapping function of the json.UnmarshalInterface function
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
	return CheckRequiredStructs(vo)
}

// CheckRequiredStructs will inspect the given reflect.Value. If it contains
// a required struct, it will check it's content, if it contains a struct
// it will recursively invoke the function once more, if none of these apply
// nil will be returned.
func CheckRequiredStructs(vo reflect.Value) error {
	if vo.Kind() != reflect.Struct {
		return nil
	}
	for i := 0; i < vo.NumField(); i++ {
		vtf := vo.Field(i)
		switch vtf.Type() {
		case reflect.TypeOf(String{}):
			return checkRequiredValue(vtf)
		}
		if vtf.Kind() == reflect.Struct {
			if err := CheckRequiredStructs(vtf); err != nil {
				return err
			}
		}
	}
	return nil
}

func checkRequiredValue(vo reflect.Value) error {
	for i := 0; i < vo.NumField(); i++ {
		vtf := vo.Field(i)
		switch vtf.Kind() {
		case reflect.String:
			if vtf.Len() == 0 {
				return ErrStringEmpty
			}
		}
	}
	return nil
}
