package json

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"github.com/Pungyeon/required/pkg/lexer"
	"github.com/Pungyeon/required/pkg/required"
	"github.com/Pungyeon/required/pkg/structtag"
	"github.com/Pungyeon/required/pkg/token"
)

func Parse(l *lexer.Lexer, v interface{}) error {
	val := getReflectValue(v)
	if val.Kind() == reflect.Interface {
		p := &parser{lexer: l}
		obj, err := p.parse(val)
		if err != nil {
			return err
		}
		val.Set(obj)
		return nil
	}
	p := &parser{lexer: l}
	if err := p.next(); err != nil {
		return err
	}
	if err := p.decode(val); err != nil {
		return err
	}
	return nil
}

func (p *parser) decode(val reflect.Value) error {
	tags, err := structtag.FromValue(val)
	if err != nil {
		return err
	}
	if tags.UnmarshalInterface {
		if val.CanAddr() {
			str := p.lexer.SkipValue()
			return val.Addr().Interface().(json.Unmarshaler).UnmarshalJSON(str)
		} else {
			return val.Interface().(json.Unmarshaler).UnmarshalJSON(p.lexer.SkipValue())
		}
	}

	if err := p._decode(val, tags); err != nil {
		return err
	}

	if tags.RequiredInterface {
		if val.CanAddr() {
			return val.Addr().Interface().(required.Required).IsValueValid()
		} else {
			return val.Interface().(required.Required).IsValueValid()
		}
	}
	return nil
}

func checkIfEOF(err error) error {
	if err == io.EOF {
		return nil
	}
	return err
}

func (p *parser) _decode(val reflect.Value, tags structtag.Tags) error {
	if val.Kind() == reflect.Interface {
		obj, err := p.parse(val)
		if err != nil {
			return err
		}
		val.Set(obj)
		return nil
	}

	if p.current.Type == token.Null {
		return nil
	}

	switch val.Kind() {
	case reflect.Ptr:
		vo := getValueOfPointer(val)
		if err := p._decode(getElemOfValue(vo), tags); err != nil {
			return err
		}
		val.Set(vo)
		return nil
	case reflect.Array, reflect.Slice:
		if err := p.decodeArray(val); err != nil {
			return err
		}
		return checkIfEOF(p.next())
	case reflect.Map:
		return p.decodeMap(val)
	case reflect.Struct:
		if err := p.decodeObject(val, tags); err != nil {
			return err
		}
		return checkIfEOF(p.next())
	case reflect.Int, reflect.Float32, reflect.Float64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.String, reflect.Bool: // what about uints?
		if err := p.current.SetValueOf(val); err != nil {
			return err
		}
		return checkIfEOF(p.next())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := token.Ttoi(p.current)
		if err != nil {
			return err
		}
		val.SetUint(uint64(n))
		return err
	}
	return token.Error(token.ErrInvalidJSON, p.current.ToString())
}

func (p *parser) decodeObject(val reflect.Value, tags structtag.Tags) error {
	if err := p.next(); err != nil {
		return checkIfEOF(err)
	}
	for {
		if p.current.Type != token.String {
			if p.current.Type == token.ClosingCurly {
				return tags.CheckRequired()
			}
			return token.Error(token.ErrInvalidJSON, fmt.Sprintf("expected object field, got: %s", p.current))
		}
		field := p.current
		if err := p.next(); err != nil {
			return checkIfEOF(err)
		}
		if p.current.Type != token.Colon {
			return token.Error(token.ErrInvalidJSON, fmt.Sprintf("expected colon token: %s", p.current))
		}
		if err := p.next(); err != nil {
			return checkIfEOF(err)
		}
		tag := tags.Tags[field.ToString()]
		if err := p.decode(val.Field(tag.FieldIndex)); err != nil {
			return err
		}

		if p.current.Type == token.Null {
			if err := p.next(); err != nil {
				return checkIfEOF(err)
			}
		} else {
			tags.Set(tag)
		}
		if p.current.Type == token.Comma {
			if err := p.next(); err != nil {
				return checkIfEOF(err)
			}
		} else {
			if p.current.Type == token.ClosingCurly {
				return tags.CheckRequired()
			} else {
				return token.Error(token.ErrInvalidJSON, fmt.Sprintf("expected closing curly or comma: (%s) -> %s", p.current, p.lexer.Previous()))
			}
		}
	}
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
	tags, err := structtag.FromValue(arr.Index(0))
	if err != nil {
		return err
	}

	var i int
	for {
		if err := p.next(); err != nil {
			return err
		}
		switch p.current.Type {
		case token.Comma:
			continue // skip commas
		case token.ClosingBrace:
			arr.Set(arr.Slice(0, i))
			return nil
		case token.OpenCurly:
			arr.Set(grow(arr, i))
			if arr.Index(i).Kind() == reflect.Map {
				if err := p.decodeMap(arr.Index(i)); err != nil {
					return err
				}
			} else {
				if err := p.decodeObject(arr.Index(i), tags); err != nil {
					return err
				}
			}
			i++
			if p.current.Type == token.ClosingBrace {
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
				val, err := p.current.AsValue(arr.Type().Elem())
				if err != nil {
					return err
				}
				arr.Index(i).Set(val)
			} else {
				if err := p.current.SetValueOf(arr.Index(i)); err != nil {
					return err
				}
			}
			i++
		}
	}
}

type parser struct {
	lexer *lexer.Lexer
	obj   reflect.Value
	current token.Token
	previous token.Token
}

func (p *parser) eof() bool {
	return p.lexer.EOF()
}

func (p *parser) next() error {
	var err error
	p.previous = p.current
	p.current, err = p.lexer.Next()
	return err
}

func (p *parser) parse(vo reflect.Value) (reflect.Value, error) {
	for {
		if err := p.next(); err != nil {
			return vo, checkIfEOF(err)
		}
		switch p.current.Type {
		case token.OpenBrace:
			return p.parseArray(determineArrayType(vo))
		case token.OpenCurly:
			return p.parseObject(vo)
		case token.Null:
			return vo, nil
		default:
			if vo.Kind() == reflect.Interface {
				return p.current.AsValue(vo.Type())
			}
			return vo, p.current.SetValueOf(vo)
		}
	}
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
	var (
		arr = reflect.MakeSlice(sliceType, 3, 3)
		val = reflect.New(sliceType.Elem()).Elem()
		i int
	)
	for {
		if err := p.next(); err != nil {
			return arr.Slice(0, i), checkIfEOF(err)
		}
		switch p.current.Type {
		case token.Comma:
			continue // skip commas
		case token.ClosingBrace:
			return arr.Slice(0, i), nil
		case token.OpenCurly:
			inner, err := p.parseMap(val.Type())
			if err != nil {
				return arr, err
			}
			arr = insertAt(arr, i, inner)
			i++
			if p.current.Type == token.ClosingBrace {
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
				val, err = p.current.AsValue(sliceType.Elem())
				if err != nil {
					return val, nil
				}
			} else {
				if err := p.current.SetValueOf(val); err != nil {
					return val, nil
				}
			}
			arr = insertAt(arr, i, val)
			i++
		}
	}
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
	for {
		if err := p.next(); err != nil {
			return checkIfEOF(err)
		}

		if p.current.Type == token.Colon {
			tag := tags.Tags[p.previous.ToString()]
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
		}
		if p.eof() || p.current.Type == token.ClosingCurly {
			if err := p.next(); err != nil {
				return checkIfEOF(err)
			}
			return tags.CheckRequired()
		}
	}
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
		val   = reflect.New(vmap.Type().Elem()).Elem()
		field = reflect.New(token.ReflectTypeString).Elem()
		err   error
	)

	if vmap.IsNil() {
		vmap.Set(reflect.MakeMap(reflect.MapOf(token.ReflectTypeString, vmap.Type().Elem())))
	}

	for {
		if err := p.next(); err != nil {
			return checkIfEOF(err)
		}
		if p.current.Type == token.ClosingCurly {
			return checkIfEOF(p.next())
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
}

func (p *parser) parseMap(valueType reflect.Type) (reflect.Value, error) {
	var (
		val   = reflect.New(valueType).Elem()
		field = reflect.New(token.ReflectTypeString).Elem()
		err   error
	)

	vmap := reflect.MakeMap(reflect.MapOf(token.ReflectTypeString, valueType))

	for {
		if err := p.next(); err != nil {
			return vmap, checkIfEOF(err)
		}
		if p.current.Type == token.ClosingCurly {
			if err := p.next(); err != nil {
				return vmap, checkIfEOF(err)
			}
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
}

func (p *parser) setField(val reflect.Value) error {
	for {
		if err := p.next(); err != nil {
			return checkIfEOF(err)
		}
		if p.current.Type == token.Colon {
			return p.previous.SetValueOf(val)
		}
	}
}
