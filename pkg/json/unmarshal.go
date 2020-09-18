package json

import (
	"fmt"
	"reflect"
	"strconv"
)

var (
	SPACE byte = ' '
	TAB byte = '\t'
	NEWLINE byte = '\n'
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
		fmt.Printf("(%d) %s: %s\n", i, tag, f.Name)
		tags[tag] = i
	}
	for r.Next() {
		switch r.Value() {
		case RIGHT_CURLY:
			r.Next()
			return nil
		case QUOTATION:
			fieldName := r.getFieldName()
			fmt.Println(fieldName)
			r.parseFieldValue(tags[fieldName], vo)
		case SPACE, NEWLINE, TAB:
			continue
		}
	}
	return nil
}


func (r *Reader) parseFieldValue(field int, value reflect.Value) {
	r.Seek(COLON)

	var isString bool
	var hasFloat bool

	var buf []byte
	for r.Next() {
		switch r.Value() {
		case LEFT_CURLY:
			v := value.Field(field)
			if err := r.ParseValue(v); err != nil {
				panic(err)
			}
			return
		case FULLSTOP:
			hasFloat = true
			buf = append(buf, r.Value())
		case SPACE:
			if !isString {
				continue
			}
		case QUOTATION:
			isString = true
		case COMMA, RIGHT_CURLY:
			goto PARSE
		case TAB, NEWLINE:
			continue
		default:
			buf = append(buf, r.Value())
		}
	}
	PARSE:

	if isString {
		value.Field(field).SetString(string(buf))
	} else {
		if hasFloat {
			val, err := strconv.ParseFloat(string(buf), 64)
			if err != nil {
				panic(err)
			}
			value.Field(field).SetFloat(val)
		} else {
			val, err := strconv.ParseInt(string(buf), 10, 64)
			if err != nil {
				panic(err)
			}
			value.Field(field).SetInt(val)
		}
	}
}

func (r *Reader) getFieldName() string {
	r.Next()
	return r.StringUntil(QUOTATION)
}


