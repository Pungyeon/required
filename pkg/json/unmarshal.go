package json

import (
	"fmt"
	"reflect"
	"strconv"
)

func getReflectValue(v interface{}) reflect.Value {
	vo := reflect.ValueOf(v)
	for vo.Kind() == reflect.Ptr {
		vo = vo.Elem()
	}
	return vo
}

func getFieldTags(vo reflect.Value) map[string]int {
	if vo.Kind() != reflect.Struct {
		return map[string]int{}
	}
	tags := make(map[string]int)
	for i := 0; i < vo.NumField(); i++ {
		f := vo.Type().Field(i)
		tag := f.Tag.Get("json")
		// TODO: if there is no tag, then assume the default tag
		tags[tag] = i
	}
	return tags
}

type ObjectType int

const (
	Unknown = 0
	String  = 1
	Integer = 2
	Float   = 3
	Slice   = 4
	Obj     = 5
)

type Object struct {
	Value string
	Type  ObjectType
}

func (obj *Object) add(token Token) {
	if token.Type == StringToken {
		obj.Type = String
	}
	if token.Type == FullStopToken {
		obj.Type = Float
	}
	obj.Value += token.Value
}

func Parse(tokens Tokens, v interface{}) error {
	fmt.Println(tokens.Join(";"))
	vo := getReflectValue(v)
	p := &parser{
		index:  -1,
		tokens: tokens,
	}

	return p.parse(vo)
}

func (p *parser) parse(vo reflect.Value) error {
	p.obj = vo
	p.tags = getFieldTags(vo)
	fmt.Println(p.tags)

	for p.next() {
		if p.current().Type == StringToken {
			fmt.Println(p.current())
		}
		if p.current().Type == OpenBraceToken {
			arr, err := p.parseArray(vo.Type())
			if err != nil {
				return err
			}
			vo.Set(arr)
			return nil
		}
		if p.current().Value == ":" {
			if err := p.setValueOnField(p.previous().Value); err != nil {
				return err
			}
		}
		if p.current().Type == ClosingCurlyToken {
			p.next()
			return nil
		}
	}
	return nil
}

type parser struct {
	tokens Tokens
	index  int
	tags   map[string]int
	obj    reflect.Value
}

func (p *parser) previous() Token {
	return p.tokens[p.index-1]
}

func (p *parser) current() Token {
	return p.tokens[p.index]
}

func (p *parser) next() bool {
	p.index++
	return p.index < len(p.tokens)
}

func (p *parser) setValueOnField(field string) error {
	for p.next() {
		switch p.current().Type {
		case OpenBraceToken:
			fmt.Println(p.tags)
			obj := p.obj.Field(p.tags[field])
			fmt.Printf("field: %s, obj: %s\n", field, obj.Type())
			arr, err := p.parseArray(obj.Type())
			if err != nil {
				return err
			}
			obj.Set(arr)
			return nil
			//return p.setArray(field)
		case OpenCurlyToken:
			fmt.Println("before:", p.tags)
			fmt.Println("after:", p.tags)
			return p.setInnerObject(field)
		default:
			return p.setPrimitive(field)
		}
	}
	return fmt.Errorf("could not parse value following: %v", field)
}

func (p *parser) parseArray(sliceType reflect.Type) (reflect.Value, error) {
	var slice []Object
	obj := &Object{Type: Integer}
	for p.next() {
		switch p.current().Type {
		case CommaToken, ClosingCurlyToken, ClosingBraceToken:
			slice = append(slice, *obj)
			obj = &Object{Type: Integer}
			if p.current().Type == ClosingBraceToken {
				goto SET_ARRAY
			}
		case OpenCurlyToken:
			// dinner object
		case OpenBraceToken:
			// inner array
		default:
			obj.add(p.current())
		}
	}
SET_ARRAY:
	fmt.Println(sliceType)
	arr := reflect.MakeSlice(sliceType, len(slice), len(slice))
	for i, obj := range slice {
		switch obj.Type {
		case String:
			arr.Index(i).SetString(obj.Value)
		case Integer:
			v, err := strconv.ParseInt(obj.Value, 10, 64)
			if err != nil {
				return arr, err
			}
			arr.Index(i).SetInt(v)
		case Float:
			v, err := strconv.ParseFloat(obj.Value, 64)
			if err != nil {
				return arr, err
			}
			arr.Index(i).SetFloat(v)
		}
	}

	return arr, nil
}

func (p *parser) setInnerObject(field string) error {
	inner := &parser{
		index:  p.index,
		tokens: p.tokens,
	}
	obj := p.obj.Field(p.tags[field])
	if err := inner.parse(obj); err != nil {
		return err
	}
	p.index = inner.index
	fmt.Println("returning setInnerObject")
	return nil
}

func (p *parser) setPrimitive(field string) error {
	str := p.current().Value
	for p.next() {
		if p.current().Type == CommaToken || p.current().Type == ClosingCurlyToken {
			fmt.Println("setting", field, str)
			setField(p.obj, p.tags[field], str)
			return nil
		} else {
			str += p.current().Value
		}
	}
	return nil
}

func getSetIndexFn(t reflect.Type) func(int, reflect.Value, string) {
	switch t {
	case reflect.TypeOf([]int{}):
		return func(i int, arr reflect.Value, element string) {
			val, err := strconv.ParseInt(element, 10, 64)
			if err != nil {
				panic(err)
			}
			arr.Index(i).SetInt(val)
		}
	case reflect.TypeOf([]string{}):
		return func(i int, arr reflect.Value, element string) {
			arr.Index(i).SetString(element)
		}
	default:
		return func(i int, arr reflect.Value, element string) {
			arr.Index(i).SetString(element)
		}
	}
}

func setField(object reflect.Value, field int, value string) {
	t := object.Field(field).Kind()
	switch t {

	case reflect.String:
		object.Field(field).SetString(value)
	case reflect.Float64:
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			panic(err)
		}
		object.Field(field).SetFloat(val)
	case reflect.Int, reflect.Int64:
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			panic(err)
		}
		object.Field(field).SetInt(val)
	default:
		fmt.Printf("could not set field: %s (%s) as %v\n", object.Type().Field(field).Name, t, value)
	}
}
