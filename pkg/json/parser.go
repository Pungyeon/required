package json

import (
	"errors"
	"reflect"

	"github.com/Pungyeon/json-validation/pkg/token"

	"github.com/Pungyeon/json-validation/pkg/structtag"

	"github.com/Pungyeon/json-validation/pkg/required"
)

func Parse(tokens token.Tokens, v interface{}) error {
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
	tokens token.Tokens
	index  int
	obj    reflect.Value
}

func (p *parser) previous() token.Token {
	return p.tokens[p.index-1]
}

func (p *parser) current() token.Token {
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
		case token.OpenBrace:
			return p.parseArray(determineArrayType(vo))
		case token.OpenCurly:
			return p.parseObject(vo)
		case token.Null:
			return vo, nil
		default:
			return p.current().AsValue(vo.Type())
		}
	}
	return reflect.New(token.ReflectTypeString), nil // shouldn't this be an error?
}

func (p *parser) parseObject(vo reflect.Value) (reflect.Value, error) {
	kind, _type := determineObjectType(vo)
	if kind == reflect.Map {
		return p.parseMap(_type)
	}
	return p.parseStructureWithCopy(vo)
}

func (p *parser) parseStructureWithCopy(vo reflect.Value) (reflect.Value, error) {
	index, err := p.copy().parseStructure(vo)
	if err != nil {
		return vo, err
	}
	p.index = index
	return vo, nil
}

func determineObjectType(vo reflect.Value) (reflect.Kind, reflect.Type) {
	if vo.Kind() == reflect.Interface {
		return reflect.Map, token.ReflectTypeInterface
	} else if vo.Kind() == reflect.Map {
		return reflect.Map, vo.Type().Elem()
	} else {
		return vo.Kind(), vo.Type()
	}
}

func determineArrayType(vo reflect.Value) reflect.Type {
	if vo.Kind() == reflect.Interface {
		return reflect.ValueOf([]interface{}{}).Type()
	}
	return vo.Type()
}

func (p *parser) parseArray(sliceType reflect.Type) (reflect.Value, error) {
	var slice []reflect.Value
	for p.next() {
		switch p.current().Type {
		case token.Comma:
			continue // skip commas
		case token.ClosingBrace:
			return p.setArray(sliceType, slice)
		case token.OpenCurly:
			obj := reflect.New(sliceType.Elem()).Elem()
			index, err := p.copy().parseStructure(obj)
			if err != nil {
				return obj, nil
			}
			p.index = index
			slice = append(slice, obj)
			if p.current().Type == token.ClosingBrace {
				return p.setArray(sliceType, slice)
			}
		case token.OpenBrace:
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
	index, err := p.copy().parseStructure(getElemOfValue(ptr))
	if err != nil {
		return index, err
	}
	vo.Set(ptr)
	return index, err
}

func (p *parser) parseStructure(vo reflect.Value) (int, error) {
	tags, err := structtag.FromValue(vo)
	if err != nil {
		return -1, err
	}
	if vo.Kind() == reflect.Ptr {
		return p.parsePointerObject(vo)
	}
	for p.next() {
		if p.current().Value == ":" {
			tag := tags[p.previous().Value]
			obj := vo.Field(tag.FieldIndex)
			val, err := p.parse(obj)
			if err != nil {
				return p.index, err
			}
			tags.Set(tag) // TODO : Make sure to not set this, if the token is a NullToken
			obj.Set(val)
			if req, ok := obj.Interface().(required.Required); ok {
				if err := req.IsValueValid(); err != nil {
					return p.index, err
				}
			}
		}
		if p.eof() || p.current().Type == token.ClosingCurly {
			p.next()
			return p.index, tags.CheckRequired()
		}
	}
	return p.index, tags.CheckRequired()
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
	vmap := reflect.MakeMap(reflect.MapOf(token.ReflectTypeString, valueType))
	for p.next() {
		if p.current().Type == token.ClosingCurly {
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
		if p.current().Type == token.Colon {
			val := reflect.New(token.ReflectTypeString).Elem()
			val.SetString(p.previous().Value)
			return val, nil
		}
	}
	return reflect.New(token.ReflectTypeString).Elem(), errors.New("could not parse field")
}
