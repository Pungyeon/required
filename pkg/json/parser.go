package json

import (
	"errors"
	"reflect"
)

var (
	reflectTypeString    = reflect.TypeOf("")
	reflectTypeInteger   = reflect.TypeOf(1)
	reflectTypeFloat     = reflect.TypeOf(3.2)
	reflectTypeInterface = reflect.ValueOf(map[string]interface{}{}).Type().Elem()
	reflectTypeBool      = reflect.TypeOf(true)
)

func Parse(tokens Tokens, v interface{}) error {
	vo := getReflectValue(v)
	p := &parser{
		index:  -1,
		tokens: tokens,
	}
	obj, err := p.parse(vo)
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

func (p *parser) parse(vo reflect.Value) (reflect.Value, error) {
	for p.next() {
		switch p.current().Type {
		case OpenBraceToken:
			if vo.Type() == nil {
				return p.parseArray(reflect.ValueOf([]interface{}{}).Type())
			} else {
				return p.parseArray(vo.Type())
			}
		case OpenCurlyToken:
			if vo.Kind() == reflect.Interface { // assuming that it's an interface type
				return p.parseMap(reflectTypeInterface)
			} else if vo.Kind() == reflect.Map {
				obj, err := p.parseMap(vo.Type().Elem())
				if err != nil {
					return obj, err
				}
				return obj, nil
			} else {
				index, err := p.copy().parseObject(vo)
				if err != nil {
					return vo, err
				}
				p.index = index
				return vo, nil
			}
		default:
			return p.current().AsValue(vo.Type())
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
			val, err := p.current().AsValue(sliceType.Elem())
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

func getValueOfPointer(vo reflect.Value) reflect.Value {
	ptr := reflect.New(vo.Type())
	p2 := ptr.Elem()
	ptr.Elem().Set(reflect.New(p2.Type().Elem()))
	return reflect.ValueOf(ptr.Elem().Interface())
}

func (p *parser) parsePointerObject(vo reflect.Value) (int, error) {
	ptr := getValueOfPointer(vo)
	index, err := p.copy().parseObject(getElemOfValue(ptr))
	if err != nil {
		panic(err)
	}
	vo.Set(ptr)
	return index, err
}

func (p *parser) parseObject(vo reflect.Value) (int, error) {
	p.tags = getFieldTags(vo)
	if vo.Kind() == reflect.Ptr {
		return p.parsePointerObject(vo)
	}

	for p.next() {
		if p.current().Value == ":" {
			obj := vo.Field(p.tags[p.previous().Value])
			val, err := p.parse(obj)
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
		val, err := p.parse(reflect.New(valueType).Elem())
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
