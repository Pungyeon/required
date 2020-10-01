package json

import (
	"fmt"
	"reflect"
	"strconv"
)

type ObjectType int

const (
	Integer   = 0
	String    = 1
	Float     = 3
	Slice     = 4
	Obj       = 5
	Interface = 6
)

var (
	i64                  int64 = 1
	reflectTypeString          = reflect.TypeOf("")
	reflectTypeInteger         = reflect.TypeOf(i64)
	reflectTypeFloat           = reflect.TypeOf(3.2)
	reflectTypeInterface       = reflect.ValueOf(map[string]interface{}{}).Type().Elem()
)

type Object struct {
	Value interface{}
	Type  ObjectType
}

func (obj *Object) AsValue() (reflect.Value, error) {
	switch obj.Type {
	case String:
		refval := reflect.New(reflectTypeString).Elem()
		refval.SetString(obj.Value.(string))
		return refval, nil
	case Integer:
		refval := reflect.New(reflectTypeInteger).Elem()
		val, err := strconv.ParseInt(obj.Value.(string), 10, 64)
		if err != nil {
			return reflect.ValueOf(nil), err
		}
		refval.SetInt(val)
		return refval, nil
	case Float:
		refval := reflect.New(reflectTypeFloat).Elem()
		val, err := strconv.ParseFloat(obj.Value.(string), 64)
		if err != nil {
			return reflect.ValueOf(nil), err
		}
		refval.SetFloat(val)
		return refval, nil
	default:
		panic(fmt.Sprintln("could not determine object type:", obj.Type))
	}
}

func (obj *Object) add(token Token) {
	switch token.Type {
	case StringToken:
		if obj.Type != Interface {
			obj.Type = String
		}
	case FullStopToken:
		obj.Type = Float
	}
	if obj.Value == nil {
		obj.Value = token.Value
	} else {
		obj.Value = obj.Value.(string) + token.Value
	}
}
