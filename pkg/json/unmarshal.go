package json

import (
	"fmt"
	"reflect"
)

var (
	SPACE byte = ' '
	TAB byte = '\t'
	NEWLINE byte = '\n'
	LEFT_CURLY byte = '{'
	RIGHT_CURLY byte = '}'
	QUOTATION byte = '"'
	COLON byte = ':'
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
	tags := make(map[string]reflect.StructField)
	vo := reflect.ValueOf(v)
	vtf := vo
	for vo.Kind() == reflect.Ptr {
		vo = vo.Elem()
	}
	for i := 0; i < vo.NumField(); i++ {
		f := vo.Type().Field(i)
		tag := f.Tag.Get("json")
		tags[tag] = f
	}


	reader := &Reader{
		data: data,
		index: -1,
	}
	for reader.Next() {
		switch reader.Value() {
		case LEFT_CURLY:
			// new object
			continue
		case RIGHT_CURLY:
			// end object
			continue
		case QUOTATION:
			// get field name
			fieldName := reader.StringUntil(QUOTATION)

			reader.Seek(COLON)
			reader.Next() // this will crash the application at some point :sob:

			fmt.Println(fieldName)
			// get field value
			reader.Seek(QUOTATION)

			fieldValue := reader.StringUntil(QUOTATION)
			fmt.Println(fieldValue)

			reflect.Indirect(vtf).FieldByName(
				tags[fieldName].Name).SetString(fieldValue)
		case SPACE, NEWLINE, TAB:
			continue
		}
		fmt.Println(string(reader.Value()))
	}
	return nil
}

func getValue(vo reflect.Value) reflect.Value {
	for vo.Kind() == reflect.Ptr {
		vo = vo.Elem()
	}
	return vo
}

