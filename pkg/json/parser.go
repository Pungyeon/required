package json

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func Parse(tokens Tokens, v interface{}) error {
	vo := getReflectValue(v)
	p := &parser{
		index:  -1,
		tokens: tokens,
	}

	//return p.parse(vo)
	obj, err := p._parse(vo.Type())
	if err != nil {
		return err
	}
	vo.Set(obj)
	return nil
}

type parser struct {
	tokens Tokens
	index  int
	tags   map[string]int
	obj    reflect.Value
}

func (p *parser) _parse(vo reflect.Type) (reflect.Value, error) {
	for p.next() {
		switch p.current().Type {
		case OpenBraceToken:
			if vo == nil {
				return p.parseArray(reflect.ValueOf([]interface{}{}).Type())
			} else {
				return p.parseArray(vo)
			}
		case OpenCurlyToken:
			if vo == nil { // assuming that it's an interface type
				return p.parseMap(reflectTypeInterface)
			} else {
				obj := reflect.New(vo).Elem()
				if err := p.parseObject(obj); err != nil {
					return obj, err
				}
				return obj, nil
			}
		default:
			return p.parsePrimitive()
		}
	}
	return reflect.New(reflectTypeString), nil
}

func (p *parser) parse(vo reflect.Value) error {
	p.obj = getElemOfValue(vo)
	p.tags = getFieldTags(vo)

	for p.next() {
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
		if p.eof() || p.current().Type == ClosingCurlyToken {
			p.next()
			return nil
		}
	}
	return nil
}

func (p *parser) previous() Token {
	return p.tokens[p.index-1]
}

func (p *parser) current() Token {
	return p.tokens[p.index]
}

func (p *parser) eof() bool {
	return p.index >= len(p.tokens)
}

func (p *parser) next() bool {
	p.index++
	return p.index < len(p.tokens)
}

func (p *parser) peekNext() (Token, bool) {
	if p.index < len(p.tokens)-1 {
		return p.tokens[p.index+1], true
	}
	return Token{}, false
}

func (p *parser) setValueOnField(field string) error {
	for p.next() {
		switch p.current().Type {
		case OpenBraceToken:
			obj := p.obj.Field(p.tags[field])
			arr, err := p.parseArray(obj.Type())
			if err != nil {
				return err
			}
			obj.Set(arr)
			return nil
		case OpenCurlyToken:
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
			if obj.Value != nil {
				slice = append(slice, *obj)
				obj = &Object{Type: Integer}
			}
			if p.current().Type == ClosingBraceToken {
				return p.setArray(sliceType, slice)
			}
		case OpenCurlyToken:
			obj := reflect.New(sliceType.Elem()).Elem()
			if err := p.parseObject(obj); err != nil {
				return obj, nil
			}
			slice = append(slice, Object{Type: Obj, Value: obj})
			if p.current().Type == ClosingBraceToken {
				return p.setArray(sliceType, slice)
			}
		case OpenBraceToken:
			inner, err := p.parseArray(sliceType.Elem())
			if err != nil {
				return inner, err
			}
			slice = append(slice, Object{Value: inner, Type: Slice})
		default:
			obj.add(p.current())
		}
	}
	return p.setArray(sliceType, slice)
}

func (p *parser) setArray(sliceType reflect.Type, slice []Object) (reflect.Value, error) {
	arr := reflect.MakeSlice(sliceType, len(slice), len(slice))
	for i, obj := range slice {
		if sliceType == reflect.TypeOf([]interface{}{}) {
			val, err := obj.AsValue()
			if err != nil {
				return val, err
			}
			arr.Index(i).Set(val)
			continue
		}
		switch obj.Type {
		case String:
			arr.Index(i).SetString(obj.Value.(string))
		case Integer:
			v, err := strconv.ParseInt(obj.Value.(string), 10, 64)
			if err != nil {
				return arr, err
			}
			arr.Index(i).SetInt(v)
		case Float:
			v, err := strconv.ParseFloat(obj.Value.(string), 64)
			if err != nil {
				return arr, err
			}
			arr.Index(i).SetFloat(v)
		case Slice:
			arr.Index(i).Set(obj.Value.(reflect.Value))
		case Obj:
			arr.Index(i).Set(obj.Value.(reflect.Value))
		}
	}

	return arr, nil
}

func (p *parser) parseObject(obj reflect.Value) error {
	inner := &parser{
		index:  p.index,
		tokens: p.tokens,
	}
	if err := inner.parse(obj); err != nil {
		return err
	}
	p.index = inner.index
	return nil
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
	return nil
}

func (p *parser) parsePrimitive() (reflect.Value, error) {
	obj := &Object{}
	obj.add(p.current())
	for p.next() {
		if p.current().Type == CommaToken || p.current().Type == ClosingCurlyToken {
			return obj.AsValue()
		} else {
			obj.add(p.current())
		}
	}
	return obj.AsValue()
}

func (p *parser) setPrimitive(field string) error {
	str := p.current().Value
	for p.next() {
		if p.current().Type == CommaToken || p.current().Type == ClosingCurlyToken {
			p.setField(field, str)
			return nil
		} else {
			str += p.current().Value
		}
	}
	return nil
}

func getReflectValue(v interface{}) reflect.Value {
	return getElemOfValue(reflect.ValueOf(v))
}

func getElemOfValue(vo reflect.Value) reflect.Value {
	for vo.Kind() == reflect.Ptr {
		vo = vo.Elem()
	}
	return vo
}

func (p *parser) setField(field string, value string) {
	if p.obj.Kind() == reflect.Map {
		setFieldOnMap(p.obj, field, value)
		return
	}
	setFieldonStruct(p.obj, p.tags[field], value)
}

func setFieldonStruct(object reflect.Value, field int, value string) {
	obj := getElemOfValue(object)
	if err := setValueOnObject(obj.Field(field), value); err != nil {
		fmt.Printf("could not set field: %s (%s) as %v\n",
			object.Type().Field(field).Name, object.Field(field).Kind(), value)
	}
}

func setValueOnObject(field reflect.Value, value string) error {
	t := field.Kind()
	switch t {
	case reflect.String:
		field.SetString(value)
	case reflect.Float64:
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			panic(err)
		}
		field.SetFloat(val)
	case reflect.Int, reflect.Int64:
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			panic(err)
		}
		field.SetInt(val)
	default:
		return fmt.Errorf("could not set field - %v: %v", field, value)
	}
	return nil
}

func (p *parser) parseMap(valueType reflect.Type) (reflect.Value, error) {
	vmap := reflect.MakeMap(reflect.MapOf(reflectTypeString, valueType))
	field, err := p.parseField()
	if err != nil {
		return vmap, err
	}
	val, err := p._parse(valueType)
	if err != nil {
		return vmap, err
	}
	vmap.SetMapIndex(field, val)
	return vmap, nil
}

func (p *parser) parseField() (reflect.Value, error) {
	for p.next() {
		if p.current().Type == ColonToken {
			val := reflect.New(reflectTypeString).Elem()
			val.SetString(p.previous().Value)
			return val, nil
		}
	}
	return reflect.New(reflectTypeString).Elem(), errors.New("could not parse field")
}

func setFieldOnMap(object reflect.Value, field string, value string) {
	key := reflect.New(reflectTypeString).Elem()
	key.SetString(field)

	val := reflect.New(object.Type().Elem()).Elem()
	if err := setValueOnObject(val, value); err != nil {
		panic(err)
	}

	if object.IsNil() {
		object.Set(
			reflect.MakeMap(
				reflect.MapOf(reflectTypeString, object.Type().Elem())))
	}
	object.SetMapIndex(key, val)
}
