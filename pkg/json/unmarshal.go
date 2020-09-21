package json

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	SPACE byte = ' '
	TAB byte = '\t'
	NEWLINE byte = '\n'
	LEFT_BRACE byte = '['
	RIGHT_BRACE byte = ']'
	LEFT_CURLY byte = '{'
	RIGHT_CURLY byte = '}'
	QUOTATION byte = '"'
	COLON byte = ':'
	COMMA byte = ','
	FULLSTOP byte = '.'
)

type Reader struct {
	data []byte
	index int
}

func (r *Reader) Next() bool {
	r.index++
	return r.index < len(r.data)
}

func (r *Reader) Value() byte {
	return r.data[r.index]
}

func (r *Reader) StringUntil(char byte) string {
	start := r.index
	for r.Next() {
		if r.Value() == char {
			return string(r.data[start:r.index])
		}
	}
	return ""
}

func (r *Reader) Seek(char byte) {
	for r.Next() {
		if r.Value() == char {
			return
		}
	}
}

func Unmarshal(data []byte, v interface{}) error {
	vo := reflect.ValueOf(v)
	for vo.Kind() == reflect.Ptr {
		vo = vo.Elem()
	}

	reader := &Reader{
		data: data,
		index: -1,
	}
	return reader.ParseValue(vo)
}

func (r *Reader) ParseValue(vo reflect.Value) error {
	tags := make(map[string]int)
	for i := 0; i < vo.NumField(); i++ {
		f := vo.Type().Field(i)
		tag := f.Tag.Get("json")
		// TODO: if there is no tag, then assume the default tag
		//fmt.Printf("(%d) %s: %s\n", i, tag, f.Name)
		tags[tag] = i
	}
	for r.Next() {
		switch r.Value() {
		case RIGHT_CURLY:
			r.Next()
			return nil
		case QUOTATION:
			fieldName := r.getFieldName()
			r.parseFieldValue(tags[fieldName], vo)
			if r.Value() == RIGHT_CURLY {
				r.Next()
				return nil
			}
		case SPACE, NEWLINE, TAB:
			continue
		}
	}
	return nil
}


func (r *Reader) parseFieldValue(field int, value reflect.Value) {
	r.Seek(COLON)
	v, ok := r.getValue(value, field)
	if ok {
		r.SetField(value, field, v)
	}
}

func (r *Reader) SetField(object reflect.Value, field int, value string) {
	t := object.Field(field).Kind()
	switch t {
	case reflect.Array, reflect.Slice:
		elements := strings.Split(value, ",")
		arr := reflect.MakeSlice(reflect.TypeOf([]int{}), len(elements), len(elements))
		for i, element := range elements {
			val, err := strconv.ParseInt(element, 10, 64)
			if err != nil {
				panic(err)
			}
			arr.Index(i).SetInt(val)
		}
		object.Field(field).Set(arr)
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

func (r *Reader) getValue(value reflect.Value, field int) (string, bool) {
	var isString bool
	var isArray bool

	var buf []byte
	for r.Next() {
		switch r.Value() {
		case LEFT_BRACE:
			isArray = true
		case RIGHT_BRACE:
			isArray = false
		case LEFT_CURLY:
			v := value.Field(field)
			if err := r.ParseValue(v); err != nil {
				panic(err)
			}
			return "", false
		case SPACE:
			if !isString {
				continue
			}
		case QUOTATION:
			isString = true
		case RIGHT_CURLY:
			return string(buf), true
		case COMMA:
			if isArray {
				buf = append(buf, r.Value())
			} else {
				return string(buf), true
			}
		case TAB, NEWLINE:
			continue
		default:
			buf = append(buf, r.Value())
		}
	}
	return "", false
}

func (r *Reader) getFieldName() string {
	r.Next()
	return r.StringUntil(QUOTATION)
}


