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
	tags := make(map[string]int)
	for i := 0; i < vo.NumField(); i++ {
		f := vo.Type().Field(i)
		tag := f.Tag.Get("json")
		// TODO: if there is no tag, then assume the default tag
		tags[tag] = i
	}
	return tags
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

	for p.next() {
		if p.current().Type == StringToken {
			fmt.Println(p.current())
		}
		if p.current().Type == ClosingCurlyToken {
			return nil
		}
		if p.current().Value == ":" {
			if err := p.setValueOnField(p.previous().Value); err != nil {
				return err
			}
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
			return nil
		case OpenCurlyToken:
			return p.setInnerObject(field)
		default:
			return p.setPrimitive(field)
		}
	}
	return fmt.Errorf("could not parse value following: %v", field)
}

func (p *parser) setArray(field string) error {
	//case reflect.Array, reflect.Slice:
	//	elements := strings.Split(value, ",")
	//	sliceType := object.Field(field).Type()
	//	arr := reflect.MakeSlice(sliceType, len(elements), len(elements))
	//	fn := getSetIndexFn(sliceType)
	//	for i, element := range elements {
	//	fn(i, arr, element)
	//	}
	//	object.Field(field).Set(arr)
	return nil
}

func (p *parser) setInnerObject(field string) error {
	fmt.Println("found a curly!")
	inner := &parser{
		index:  p.index,
		tokens: p.tokens,
	}
	obj := p.obj.Field(p.tags[field])
	fmt.Println(obj.Type())
	if err := inner.parse(obj); err != nil {
		return err
	}
	p.index = inner.index
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
