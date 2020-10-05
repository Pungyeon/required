package json

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

var (
	i64                  int64 = 1 // TODO: How do I distinguish between i64, i32 and int?
	reflectTypeString          = reflect.TypeOf("")
	reflectTypeInteger         = reflect.TypeOf(1)
	reflectTypeFloat           = reflect.TypeOf(3.2)
	reflectTypeInterface       = reflect.ValueOf(map[string]interface{}{}).Type().Elem()
	reflectTypeBool            = reflect.TypeOf(true)
)

func Parse(tokens Tokens, v interface{}) error {
	vo := getReflectValue(v)
	p := &parser{
		index:  -1,
		tokens: tokens,
	}
	obj, err := p.parse(vo.Type())
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

func (p *parser) parse(vo reflect.Type) (reflect.Value, error) {
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
			} else if vo.Kind() == reflect.Map {
				obj, err := p.parseMap(vo.Elem())
				if err != nil {
					return obj, err
				}
				return obj, nil
			} else {
				obj := reflect.New(vo).Elem()
				index, err := p.copy().parseObject(obj)
				if err != nil {
					return obj, err
				}
				p.index = index
				return obj, nil
			}
		default:
			return p.current().ToValue()
		}
	}
	return reflect.New(reflectTypeString), nil
}

func (p *parser) parseArray(sliceType reflect.Type) (reflect.Value, error) {
	var slice []reflect.Value
	for p.next() {
		switch p.current().Type {
		case CommaToken:
			// do nothing
			continue
		case ClosingBraceToken:
			return p.setArray(sliceType, slice)
		case OpenCurlyToken:
			obj := reflect.New(sliceType.Elem()).Elem()
			index, err := p.copy().parseObject(obj)
			if err != nil {
				return obj, nil
			}
			p.index = index
			slice = append(slice, obj)
			if p.current().Type == ClosingBraceToken {
				return p.setArray(sliceType, slice)
			}
		case OpenBraceToken:
			inner, err := p.parseArray(sliceType.Elem())
			if err != nil {
				return inner, err
			}
			slice = append(slice, inner)
		default:
			val, err := p.current().ToValue()
			if err != nil {
				return val, err
			}
			slice = append(slice, val)
		}
	}
	return p.setArray(sliceType, slice)
}

func (p *parser) setArray(sliceType reflect.Type, slice []reflect.Value) (reflect.Value, error) {
	arr := reflect.MakeSlice(sliceType, len(slice), len(slice))
	for i, val := range slice {
		arr.Index(i).Set(val)
	}

	return arr, nil
}

func (p *parser) copy() *parser {
	return &parser{
		index:  p.index,
		tokens: p.tokens,
	}
}

func (p *parser) parseObject(vo reflect.Value) (int, error) {
	p.obj = getElemOfValue(vo)
	p.tags = getFieldTags(vo)

	for p.next() {
		if p.current().Value == ":" {
			obj := p.obj.Field(p.tags[p.previous().Value])
			val, err := p.parse(obj.Type())
			if err != nil {
				panic(err)
			}
			obj.Set(val)
		}
		if p.eof() || p.current().Type == ClosingCurlyToken {
			p.next()
			return p.index, nil
		}
	}
	return p.index, nil
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
	setFieldOnStruct(p.obj, p.tags[field], value)
}

func setFieldOnStruct(object reflect.Value, field int, value string) {
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
	for p.next() {
		if p.current().Type == ClosingCurlyToken {
			p.next()
			break
		}
		field, err := p.parseField()
		if err != nil {
			return vmap, err
		}
		val, err := p.parse(valueType)
		if err != nil {
			return vmap, err
		}
		vmap.SetMapIndex(field, val)
	}
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
	val := reflect.New(object.Type().Elem()).Elem()
	if err := setValueOnObject(val, value); err != nil {
		panic(err)
	}

	if object.IsNil() {
		object.Set(newReflectMap(object))
	}

	object.SetMapIndex(
		newMapKey(field),
		val)
}

func newReflectMap(object reflect.Value) reflect.Value {
	return reflect.MakeMap(
		reflect.MapOf(reflectTypeString, object.Type().Elem()))
}

func newMapKey(field string) reflect.Value {
	key := reflect.New(reflectTypeString).Elem()
	key.SetString(field)
	return key
}
