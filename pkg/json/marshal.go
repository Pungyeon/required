package json

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func Marshal(v interface{}) ([]byte, error) {
	return marshal(v)
}

func marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := _marshal(reflect.ValueOf(v), &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

var (
	quote = '"'
	colon = ':'
)

const (
	TRUE  = `true`
	FALSE = `false`
)

var scratch [64]byte

func _marshal(val reflect.Value, buf *bytes.Buffer) error {
	switch val.Kind() {
	case reflect.Float64, reflect.Float32:
		// The standard library uses a []byte array and AppendFloat
		// see encode.go:573 -> func (bits floatEncoder) encode(e *encodeState, v reflect.Value, opts encOpts)
		// I'm not exactly sure why this saves an allocation, but it does :shrug:
		f := val.Float()
		b := scratch[:0]
		b = strconv.AppendFloat(b, f, 'f', -1, 64)
		buf.Write(b)
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		buf.WriteString(strconv.FormatInt(val.Int(), 10))
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		buf.WriteString(strconv.FormatUint(val.Uint(), 10))
		return nil
	case reflect.Bool:
		if val.Bool() {
			buf.WriteString(TRUE)
		} else {
			buf.WriteString(FALSE)
		}
		return nil
	case reflect.String:
		buf.WriteRune(quote)
		buf.WriteString(val.String())
		buf.WriteRune(quote)
		return nil
	case reflect.Struct:
		return marshalStruct(val, buf)

	case reflect.Ptr:
		if val.IsNil() {
			buf.WriteString("null")
			return nil
		}
		return _marshal(val.Elem(), buf)
	case reflect.Interface:
		if val.IsNil() {
			buf.WriteString("null")
			return nil
		}
		return _marshal(val.Elem(), buf)
	case reflect.Map:
		if val.IsNil() {
			buf.WriteString("null")
			return nil
		}
		return marshalMap(val, buf)
	case reflect.Array, reflect.Slice:
		if val.IsNil() {
			buf.WriteString("null")
			return nil
		}
		return marshalArray(val, buf)
	}

	return errUnsupportedType{val: val}
}

var ErrUnsupportedType = errors.New("(required::json) unsupported type")

type errUnsupportedType struct {
	val reflect.Value
}

func (err errUnsupportedType) Unwrap() error {
	return ErrUnsupportedType
}

func (err errUnsupportedType) Error() string {
	if err.val.IsValid() {
		return fmt.Sprintf("%v: (kind: %v) (type: %v)",
			ErrUnsupportedType, err.val.Kind(), err.val.Type())
	}
	return fmt.Sprintf("%v: (kind: %v)",
		ErrUnsupportedType, err.val.Kind())
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

func marshalMap(val reflect.Value, buf *bytes.Buffer) error {
	buf.WriteString("{")
	kv := val.MapRange()

	hasNext := kv.Next()
	for hasNext {
		if err := marshalMapField(kv.Key(), buf); err != nil {
			return err
		}
		buf.WriteRune(colon)

		if err := _marshal(kv.Value(), buf); err != nil {
			return err
		}
		hasNext = kv.Next()
		if hasNext {
			buf.WriteByte(',')
		}
	}

	buf.WriteString("}")
	return nil
}

func marshalMapField(val reflect.Value, buf *bytes.Buffer) error {
	switch val.Kind() {
	case reflect.Float64, reflect.Float32:
		// The standard library uses a []byte array and AppendFloat
		// see encode.go:573 -> func (bits floatEncoder) encode(e *encodeState, v reflect.Value, opts encOpts)
		// I'm not exactly sure why this saves an allocation, but it does :shrug:
		f := val.Float()
		b := scratch[:0]
		strconv.AppendFloat(b, f, 'f', -1, 64)
		buf.WriteRune(quote)
		buf.WriteString(string(b))
		buf.WriteRune(quote)
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		buf.WriteRune(quote)
		buf.WriteString(strconv.FormatInt(val.Int(), 10))
		buf.WriteRune(quote)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		buf.WriteRune(quote)
		buf.WriteString(strconv.FormatUint(val.Uint(), 10))
		buf.WriteRune(quote)
		return nil
	case reflect.String:
		buf.WriteRune(quote)
		buf.WriteString(val.String())
		buf.WriteRune(quote)
		return nil
	}
	return fmt.Errorf("unsupported map key: %v %v", val.Kind(), val.Type())
}

func marshalStruct(val reflect.Value, buf *bytes.Buffer) error {
	buf.WriteString("{")
	tags, err := GetJSONFieldName(val)
	if err != nil {
		return err
	}
	var i int
	for tags[i].private { // TODO : @pungyeon does this check really make sure that it's private? :|
		i++
	}
	for i < val.NumField() {
		buf.WriteString(tags[i].name)
		buf.WriteRune(colon)
		if err := _marshal(val.Field(i), buf); err != nil {
			return err
		}
		i++
		for i < val.NumField() && tags[i].private {
			i++
		}
		if i < val.NumField()-1 {
			buf.WriteByte(',')
		}
	}
	buf.WriteString("}")
	return nil
}

var fieldCache = make(map[reflect.Type][]field)

type field struct {
	private     bool
	name        string
	required    bool
	omitifempty bool
}

var diff uint8 = 'a' - 'A'

func GetJSONFieldName(val reflect.Value) ([]field, error) {
	var f reflect.StructField
	tags, ok := fieldCache[val.Type()]
	if ok {
		return tags, nil
	}
	tags = make([]field, val.NumField())
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
			tags[i] = field{
				private: f.PkgPath != "",
				name:    buf.String(),
			}
		} else {
			var s, c = 0, 0
			for c < len(jsonTag) {
				if jsonTag[c] == ',' {
					if tags[i].name == "" {
						tags[i].name = `"` + jsonTag[s:c] + `"`
					} else {
						switch jsonTag[s:c] {
						case "required":
							tags[i].required = true
						case "omitifempty":
							tags[i].omitifempty = true
						default:
							return tags, fmt.Errorf("illegal json tag: %v", jsonTag)
						}
					}
					s = c
				}
				c++
			}
			tags[i].private = f.PkgPath != ""
		}
	}
	fieldCache[val.Type()] = tags
	return tags, nil
}
