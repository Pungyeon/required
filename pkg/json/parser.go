package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/Pungyeon/required/pkg/lexer"
	"github.com/Pungyeon/required/pkg/required"
	"github.com/Pungyeon/required/pkg/structtag"
	"github.com/Pungyeon/required/pkg/token"
)

/* TODO
Backtrack on architecture:
	- instead of parsing the type from the json, you should determine the type ahead of time based on the field type of the value. Of course, if it's an interface, you will have to parse the type... :shrug:
*/

func Parse(l *lexer.Lexer, v interface{}) error {
	val := getReflectValue(v)
	if val.Kind() == reflect.Interface {
		obj, err := (&parser{lexer: l}).parse(val)
		if err != nil {
			return err
		}
		val.Set(obj)
		return nil
	}
	return (&parser{lexer: l}).decode(val)
}

func (p *parser) decode(val reflect.Value) error {
	for p.next() {
		switch p.current().Type {
		case token.OpenBrace:
			return p.decodeArray(val)
		case token.OpenCurly:
			if val.Kind() == reflect.Map {
				elem, err := p.parseMap(val.Type().Elem())
				if err != nil {
					return err
				}
				val.Set(elem)
				return nil
			}
			return p.decodeObject(val)
		case token.Null:
			return nil
		default:
			if val.Kind() == reflect.Interface {
				v, err := p.current().AsValue(val.Type())
				if err != nil {
					return err
				}
				val.Set(v)
				return nil
			}
			return p.current().SetValueOf(val)
		}
	}
	return nil
}

func (p *parser) decodeField(parent reflect.Value, index int) error {
	val := parent.Field(index)

	for p.next() {
		switch p.current().Type {
		case token.OpenBrace:
			return p.decodeArray(val)
		case token.OpenCurly:
			if val.IsZero() {
				//kind, _type := determineObjectType(val)
				if val.Kind() == reflect.Ptr {
					vo := getValueOfPointer(val)
					if err := p.decodeObject(getElemOfValue(vo)); err != nil {
						return err
					}
					val.Set(vo)
					return nil
				}
				//fmt.Println(val.Type())
				if val.Kind() == reflect.Interface {
					elem, err := p.parseMap(val.Type())
					if err != nil {
						return err
					}
					val.Set(elem)
					return nil
				}
				elem, err := p.parseMap(val.Type().Elem())
				if err != nil {
					return err
				}
				val.Set(elem)
				return nil
			}
			return p.decodeObject(val)
		case token.Null:
			return nil
		default:
			if val.Kind() == reflect.Interface {
				v, err := p.current().AsValue(val.Type())
				if err != nil {
					return err
				}
				val.Set(v)
				return nil
			}
			return p.current().SetValueOf(val)
		}
	}
	return nil
}

func (p *parser) decodeObject(val reflect.Value) error {
	tags, err := structtag.FromValue(val)
	if err != nil {
		return err
	}
	for p.next() {
		if p.current().Type == token.Colon {
			tag := tags[p.previous().ToString()]
			if tags[structtag.UnmarshalInterfaceKey].Required {
				fmt.Println("hell yeah!")
				if err := val.Field(tag.FieldIndex).Interface().(json.Unmarshaler).UnmarshalJSON(p.lexer.SkipValue()); err != nil {
					return err
				}
				// TODO : @pungyeon fix this shit
				tags.Set(tag)
				continue
			}
			if err := p.decodeField(val, tag.FieldIndex); err != nil {
				return err
			}

			if tags[structtag.RequiredInterfaceKey].Required {
				if err := val.Field(tag.FieldIndex).Interface().(required.Required).IsValueValid(); err != nil {
					return err
				}
			}

			tags.Set(tag)
		}
		if p.eof() || p.current().Type == token.ClosingCurly {
			p.next()
			break
		}
	}

	return tags.CheckRequired()
}

func grow(arr reflect.Value, i int) reflect.Value {
	if arr.Len() <= i {
		grown := reflect.MakeSlice(arr.Type(), i*2, i*2)
		reflect.Copy(grown, arr)
		return grown
	}
	return arr
}

func (p *parser) decodeArray(arr reflect.Value) error {
	arr.Set(reflect.MakeSlice(arr.Type(), 3, 3))

	var i int
	for p.next() {
		switch p.current().Type {
		case token.Comma:
			continue // skip commas
		case token.ClosingBrace:
			arr.Set(arr.Slice(0, i))
			return nil
		case token.OpenCurly:
			arr.Set(grow(arr, i))
			if err := p.decodeObject(arr.Index(i)); err != nil {
				return err
			}
			i++
			if p.current().Type == token.ClosingBrace {
				arr.Set(arr.Slice(0, i))
				return nil
			}
		case token.OpenBrace:
			arr.Set(grow(arr, i))
			if err := p.decodeArray(arr.Index(i)); err != nil {
				return err
			}
			i++
		default:
			arr.Set(grow(arr, i))
			// Doing this check saves ~12 allocations per op
			if arr.Type().Elem().Kind() == reflect.Interface {
				val, err := p.current().AsValue(arr.Type().Elem())
				if err != nil {
					return err
				}
				arr.Index(i).Set(val)
			} else {
				if err := p.current().SetValueOf(arr.Index(i)); err != nil {
					return err
				}
			}
			i++
		}
	}
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
			tag := tags[p.previous().ToString()]
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
func (p *parser) decodeMap(vmap reflect.Value) error {
	var (
		val   = reflect.New(vmap.Type()).Elem()
		field = reflect.New(token.ReflectTypeString).Elem()
		err   error
	)

	if vmap.IsNil() {
		vmap.Set(reflect.MakeMap(reflect.MapOf(token.ReflectTypeString, vmap.Type())))
	}

	for p.next() {
		if p.current().Type == token.ClosingCurly {
			p.next()
			break
		}
		err = p.setField(field)
		if err != nil {
			return err
		}
		val, err = p.parse(val)
		if err != nil {
			return err
		}
		vmap.SetMapIndex(field, val)
	}
	return nil
}

func (p *parser) parseMap(valueType reflect.Type) (reflect.Value, error) {
	var (
		val   = reflect.New(valueType).Elem()
		field = reflect.New(token.ReflectTypeString).Elem()
		err   error
	)

	vmap := reflect.MakeMap(reflect.MapOf(token.ReflectTypeString, valueType))

	for p.next() {
		if p.current().Type == token.ClosingCurly {
			p.next()
			break
		}
		err = p.setField(field)
		if err != nil {
			return vmap, err
		}
		val, err = p.parse(val)
		if err != nil {
			return vmap, err
		}
		vmap.SetMapIndex(field, val)
	}
	return vmap, nil
}

func (p *parser) setField(val reflect.Value) error {
	for p.next() {
		if p.current().Type == token.Colon {
			return p.previous().SetValueOf(val)
		}
	}
	return errors.New("could not parse field")
}
