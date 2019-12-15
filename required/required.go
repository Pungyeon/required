package required

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func Unmarshal(data []byte, v interface{}) error {
	return ReturnIfError(
		json.Unmarshal(data, v),
		Required(v),
	)
}

func Required(v interface{}) error {
	vo := reflect.ValueOf(v)
	for vo.Kind() == reflect.Ptr {
		vo = vo.Elem()
	}
	return CheckRequired(vo)
}

func CheckRequired(vo reflect.Value) error {
	if vo.Kind() == reflect.Struct {
		numFields := vo.NumField()
		for i := 0; i < numFields; i++ {
			vtf := vo.Field(i)
			fmt.Println(vtf.Type().String())
			switch vtf.Type().String() {
			case "required.String":
				return checkRequiredValue(vtf)
			}
			if vtf.Kind() == reflect.Struct {
				CheckRequired(vtf)
			}
		}
	}
	return nil
}

func checkRequiredValue(vo reflect.Value) error {
	nf := vo.NumField()
	for i := 0; i < nf; i++ {
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

func ReturnIfError(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}
