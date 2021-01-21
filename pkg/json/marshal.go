package json

import (
	"bytes"
	stdjson "encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func Marshal(v interface{}) ([]byte, error) {
	return stdjson.Marshal(v)
}

func marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := _marshal(reflect.ValueOf(v), &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

var (
	quote = `"`
	colon = `:`
)

const (
	TRUE  = `true`
	FALSE = `false`
)

//type encodingFn func(reflect.Value, *bytes.Buffer) error
//
//var encodeFn [27]encodingFn
//
//func init() {
//	encodeFn = [27]encodingFn{
//		unsupported,   // Invalid
//		marshalBool,   // Bool
//		marshalInt,    // Int
//		marshalInt,    // Int8
//		marshalInt,    // Int16
//		marshalInt,    // Int32
//		marshalInt,    // Int64
//		unsupported,   // Uint
//		unsupported,   // Uint8
//		unsupported,   // Uint16
//		unsupported,   // Uint32
//		unsupported,   // Uint64
//		unsupported,   // Uintptr
//		unsupported,   // Float32
//		unsupported,   // Float64
//		unsupported,   // Complex64
//		unsupported,   // Complex128
//		unsupported,   // Array
//		unsupported,   // Chan
//		unsupported,   // Func
//		unsupported,   // Interface
//		unsupported,   // Map
//		unsupported,   // Ptr
//		unsupported,   // Slice
//		marshalString, // String
//		marshalStruct, // Struct
//		unsupported,   // UnsafePointer
//	}
//}
//
//func unsupported(val reflect.Value, _ *bytes.Buffer) error {
//	return fmt.Errorf("unsupported type: %v %v", val.Kind(), val.Type())
//}
//
//func marshalBool(val reflect.Value, buf *bytes.Buffer) error {
//	if val.Bool() {
//		buf.WriteString(TRUE)
//	} else {
//		buf.WriteString(FALSE)
//	}
//	return nil
//}
//
//func marshalInt(val reflect.Value, buf *bytes.Buffer) error {
//	buf.WriteString(strconv.FormatInt(val.Int(), 10))
//	return nil
//}
//
//func marshalString(val reflect.Value, buf *bytes.Buffer) error {
//	buf.WriteString(quote + val.String() + quote)
//	return nil
//}

func _marshal(val reflect.Value, buf *bytes.Buffer) error {
	switch val.Kind() {
	case reflect.Struct:
		return marshalStruct(val, buf)
	case reflect.Array, reflect.Slice:
		return marshalArray(val, buf)
	case reflect.Float64, reflect.Float32:
		buf.WriteString(strconv.FormatFloat(val.Float(), 'f', -1, 64))
		return nil
	case reflect.Int: // TODO : Add other Integer kinds
		buf.WriteString(strconv.FormatInt(val.Int(), 10))
		return nil
	case reflect.Bool:
		if val.Bool() {
			buf.WriteString(TRUE)
		} else {
			buf.WriteString(FALSE)
		}
		return nil
	case reflect.String:
		buf.WriteString(quote + val.String() + quote)
		return nil
	case reflect.Chan:
		// do something
	}
	return errors.New("oh dear")
}

func marshalArray(val reflect.Value, buf *bytes.Buffer) error {
	buf.WriteByte('[')
	for i := 0; i < val.Len(); i++ {
		if err := _marshal(val.Index(i), buf); err != nil {
			return err
		}
		if i < val.Len()-1 {
			buf.WriteByte(',')
		}
	}
	buf.WriteByte(']')
	return nil
}

func marshalStruct(val reflect.Value, buf *bytes.Buffer) error {
	buf.WriteString("{")
	tags, err := GetJSONFieldName(val)
	if err != nil {
		return err
	}
	for i := 0; i < val.NumField(); i++ {
		buf.WriteString(tags[i] + colon)
		if err := _marshal(val.Field(i), buf); err != nil {
			return err
		}
		if i < val.NumField()-1 {
			buf.WriteByte(',')
		}
	}
	buf.WriteString("}")
	return nil
}

var fieldCache = make(map[reflect.Type][]string)

var diff uint8 = 'a' - 'A'

func GetJSONFieldName(val reflect.Value) ([]string, error) {
	var f reflect.StructField
	tags, ok := fieldCache[val.Type()]
	if ok {
		return tags, nil
	}
	tags = make([]string, val.NumField())
	for i := 0; i < val.NumField(); i++ {
		f = val.Type().Field(i)
		jsonTag, ok := f.Tag.Lookup("json")
		if !ok {
			// use string concat instead ?
			var buf bytes.Buffer
			buf.WriteRune('"')
			for i := 0; i < len(f.Name); i++ {
				if f.Name[i] >= 'A'-1 && f.Name[i] <= 'Z' {
					if i > 0 {
						buf.WriteByte('_')
					}
					buf.WriteByte(f.Name[i] + diff)
				} else {
					buf.WriteByte(f.Name[i])
				}
			}
			buf.WriteRune('"')
			tags[i] = buf.String()
		} else {
			term := indexOf(jsonTag, '"')
			if term == -1 || term >= len(jsonTag) {
				return tags, fmt.Errorf("illegal json tag: %v", jsonTag)
			}
			tags[i] = jsonTag[jsonPrefixLen : term+1]
		}
	}
	fieldCache[val.Type()] = tags
	return tags, nil
}

const jsonPrefixLen = len(`json:"`)

func indexOf(input string, b byte) int {
	var (
		i = jsonPrefixLen
	)
	for i < len(input) {
		if input[i] == b {
			return i
		}
		i++
	}
	return -1
}
