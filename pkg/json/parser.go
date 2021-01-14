package json

import (
	"errors"
	"reflect"

	"github.com/Pungyeon/required/pkg/required"

	"github.com/Pungyeon/required/pkg/lexer"

	"github.com/Pungyeon/required/pkg/structtag"
	"github.com/Pungyeon/required/pkg/token"
)

func Parse(l *lexer.Lexer, v interface{}) error {
	vo := getReflectValue(v)
	obj, err := (&parser{lexer: l}).parse(vo)
	if err != nil {
		return err
	}
	vo.Set(obj)
	return nil
}

type parser struct {
	lexer *lexer.Lexer
	obj   reflect.Value
}

func (p *parser) previous() token.Token {
	return p.lexer.Previous()
}

func (p *parser) current() token.Token {
	return p.lexer.Current()
}

func (p *parser) eof() bool {
	return p.lexer.EOF()
}

func (p *parser) next() bool {
	return p.lexer.Next()
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
			if vo.Kind() == reflect.Interface {
				return p.current().AsValue(vo.Type())
			}
			return vo, p.current().SetValueOf(vo)
		}
	}
	return reflect.New(token.ReflectTypeString), nil // shouldn't this be an error?
}

func (p *parser) parseObject(vo reflect.Value) (reflect.Value, error) {
	kind, _type := determineObjectType(vo)
	if kind == reflect.Map {
		return p.parseMap(_type)
	}
	return vo, p.parseStructure(vo)
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

func insertAt(arr reflect.Value, i int, val reflect.Value) reflect.Value {
	if i < arr.Len() {
		arr.Index(i).Set(val)
		return arr
	}
	return reflect.Append(arr, val)
}

func (p *parser) parseArray(sliceType reflect.Type) (reflect.Value, error) {
	arr := reflect.MakeSlice(sliceType, 3, 3)
	val := reflect.New(sliceType.Elem()).Elem()
	var i int
	for p.next() {
		switch p.current().Type {
		case token.Comma:
			continue // skip commas
		case token.ClosingBrace:
			return arr.Slice(0, i), nil
		case token.OpenCurly:
			if err := p.parseStructure(val); err != nil {
				return val, nil
			}
			arr = insertAt(arr, i, val)
			i++
			if p.current().Type == token.ClosingBrace {
				return arr.Slice(0, i), nil
			}
		case token.OpenBrace:
			inner, err := p.parseArray(sliceType.Elem())
			if err != nil {
				return inner, err
			}
			arr = insertAt(arr, i, inner)
			i++
		default:
			// Doing this check saves ~12 allocations per op
			if sliceType.Elem().Kind() == reflect.Interface {
				var err error
				val, err = p.current().AsValue(sliceType.Elem())
				if err != nil {
					return val, nil
				}
			} else {
				if err := p.current().SetValueOf(val); err != nil {
					return val, nil
				}
			}
			arr = insertAt(arr, i, val)
			i++
		}
	}
	return arr.Slice(0, i), nil
}

func getValueOfPointer(vo reflect.Value) reflect.Value {
	a := reflect.New(vo.Type())
	b := a.Elem()
	a.Elem().Set(reflect.New(b.Type().Elem()))
	return reflect.ValueOf(a.Elem().Interface())
}

func (p *parser) parsePointerObject(vo reflect.Value) error {
	ptr := getValueOfPointer(vo)
	if err := p.parseStructure(getElemOfValue(ptr)); err != nil {
		return err
	}
	vo.Set(ptr)
	return nil
}

func (p *parser) parseStructure(vo reflect.Value) error {
	tags, err := structtag.FromValue(vo)
	if err != nil {
		return err
	}
	if vo.Kind() == reflect.Ptr {
		return p.parsePointerObject(vo)
	}
	for p.next() {
		if p.current().Type == token.Colon {
			tag := tags[p.previous().Value.(string)]
			obj := vo.Field(tag.FieldIndex)
			if !obj.CanSet() { // Private values may not be set.
				continue
			}
			val, err := p.parse(obj)
			if err != nil {
				return err
			}
			tags.Set(tag) // TODO : Make sure to not set this, if the token is a NullToken
			obj.Set(val)
			// TODO : This is currently 10+ allocations per op :/
			if req, ok := obj.Interface().(required.Required); ok {
				if err := req.IsValueValid(); err != nil {
					return err
				}
			}
		}
		if p.eof() || p.current().Type == token.ClosingCurly {
			p.next()
			return tags.CheckRequired()
		}
	}
	return tags.CheckRequired()
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
			val.SetString(p.previous().Value.(string))
			return val, nil
		}
	}
	return reflect.New(token.ReflectTypeString).Elem(), errors.New("could not parse field")
}
