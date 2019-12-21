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

var (
	StringType       = reflect.TypeOf(String{})
	IntType          = reflect.TypeOf(Int{})
	BoolType         = reflect.TypeOf(Bool{})
	Float32Type      = reflect.TypeOf(Float32{})
	Float64Type      = reflect.TypeOf(Float64{})
	StringSliceType  = reflect.TypeOf(StringSlice{})
	ByteSliceType    = reflect.TypeOf(ByteSlice{})
	IntSlicetype     = reflect.TypeOf(IntSlice{})
	Float32SliceType = reflect.TypeOf(Float32Slice{})
	Float64SliceType = reflect.TypeOf(Float64Slice{})
	BoolSliceType    = reflect.TypeOf(BoolSlice{})
)

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
		switch vtf.Type() {
		case StringType, IntType, BoolType,
			Float32Type, Float64Type,
			StringSliceType, ByteSliceType,
			Float32SliceType, Float64SliceType,
			IntSlicetype, BoolSliceType:
			return checkRequiredValue(vtf)
		}
		if vtf.Kind() == reflect.Struct {
			if err := CheckStructIsRequired(vtf); err != nil {
				return err
			}
		}
	}
	return nil
}

func checkRequiredValue(vo reflect.Value) error {
	// TODO : Something is really wrong here, and this implementation,
	// doens't work on embedded types for some reason ?
	// In other words, this works
	// struct {
	// 	active Bool
	// }
	// but this doesn't
	// struct {
	// 	Bool
	// }
	// for some reason the vtf.IsNil() no longer evaluates to true :|
	for i := 0; i < vo.NumField(); i++ {
		vtf := vo.Field(i)
		switch vtf.Kind() {
		case reflect.Ptr, reflect.Slice:
			if vtf.IsNil() {
				return ErrEmpty
			}
		case reflect.String:
			if vtf.Len() == 0 {
				return ErrStringEmpty
			}
		}
	}
	return nil
}
