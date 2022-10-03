package token

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

var (
	Space     byte = ' '
	Tab       byte = '\t'
	NewLine   byte = '\n'
	Quotation byte = '"'
	Escape    byte = '\\'

	ErrValueMismatch = errors.New("cannot set value of specified variable")
)

var (
	ReflectTypeString    = reflect.TypeOf("")
	ReflectTypeInteger   = reflect.TypeOf(1)
	ReflectTypeFloat     = reflect.TypeOf(3.2)
	ReflectTypeInterface = reflect.ValueOf(map[string]interface{}{}).Type().Elem()
	ReflectTypeBool      = reflect.TypeOf(true)

	ErrInvalidValue   = errors.New("invalid token value")
	ErrInvalidJSON    = errors.New("invalid json")
	ErrUnmatchedBrace = errors.New("unmatched brace found")
	ErrMissingBrace   = errors.New("missing closing brace")

	Empty = Token{}
)

type tokenErr struct {
	err     error
	details string
}

func Error(err error, details string) error {
	return tokenErr{err, details}
}

func (err tokenErr) Error() string {
	return fmt.Sprintf("%v: %v", err.err, err.details)
}

func (err tokenErr) Unwrap() error {
	return err.err
}

func Ttoi(token Token) (int64, error) {
	n, err := strconv.ParseInt(token.ToString(), 10, 64)
	if err != nil {
		return 0, Error(ErrInvalidValue, fmt.Sprintf("%v: %v", token.String(), err.Error()))
	}
	return n, nil
}

func Ttof(token Token) (float64, error) {
	n, err := strconv.ParseFloat(token.ToString(), 64)
	if err != nil {
		return 0, Error(ErrInvalidValue, fmt.Sprintf("%v: %v", token.String(), err.Error()))
	}
	return n, nil
}

type TokenType int

func (t TokenType) IsEnding() bool {
	return t == ClosingBrace || t == ClosingCurly ||
		t == ClosingBracket
}

func (t TokenType) IsOpening() bool {
	return t == OpenBrace || t == OpenCurly ||
		t == OpenBracket
}

const (
	Unknown TokenType = iota
	Integer
	Float
	String
	Null
	Colon
	Comma
	OpenBrace
	ClosingBrace
	OpenBracket
	ClosingBracket
	OpenCurly
	ClosingCurly
	FullStop
	Boolean
)

func (t TokenType) String() string {
	switch t {
	case Integer:
		return "Integer"
	case Float:
		return "Float"
	case String:
		return "String"
	case Null:
		return "Null"
	case Colon:
		return "Colon"
	case Comma:
		return "Comma"
	case OpenBrace:
		return "OpenBrace"
	case ClosingBrace:
		return "ClosingBrace"
	case OpenBracket:
		return "OpenBracket"
	case ClosingBracket:
		return "ClosingBracket"
	case OpenCurly:
		return "OpenCurly"
	case ClosingCurly:
		return "ClosingCurly"
	case FullStop:
		return "FullStop"
	case Boolean:
		return "Boolean "
	default:
		return fmt.Sprintf("UNKNOWN: %d", t)
	}
}

func init() {
	TokenTypes[':'] = Colon
	TokenTypes[','] = Comma
	TokenTypes['['] = OpenBrace
	TokenTypes[']'] = ClosingBrace
	TokenTypes['('] = OpenBracket
	TokenTypes[')'] = ClosingBracket
	TokenTypes['{'] = OpenCurly
	TokenTypes['}'] = ClosingCurly
	TokenTypes['.'] = FullStop
}

var TokenTypes = make([]TokenType, 126)

var BraceOpposites = map[byte]byte{
	'[': ']',
	']': '[',
	'(': ')',
	')': '(',
	'{': '}',
	'}': '{',
}

type Token struct {
	Value []byte
	Type  TokenType
}

func NewToken(b []byte, i int) Token {
	return Token{
		Value: b[i : i+1], // should we even allocate here?
		Type:  TokenTypes[b[i]],
	}
}

func (token Token) AsValue(vt reflect.Type) (reflect.Value, error) {
	if vt == ReflectTypeInterface || vt == nil {
		return token.ToValue()
	}

	val := reflect.New(vt).Elem()
	switch token.Type {
	case String:
		val.SetString(string(token.Value))
		return val, nil
	case Integer:
		n, err := Ttoi(token)
		if err != nil {
			return val, err
		}
		val.SetInt(n)
		return val, err
	case Float:
		f, err := strconv.ParseFloat(string(token.Value), 64)
		if err != nil {
			return val, err
		}
		val.SetFloat(f)
		return val, err
	case Boolean:
		// TODO : this doesn't actually work it will parse any value starting with 'f' and/or 't' as boolean
		// so "boolean": faulty -> false
		switch token.Value[0] {
		case 't':
			val.SetBool(true)
		case 'f':
			val.SetBool(false)
		default:
			return val, Error(ErrInvalidValue, token.String())
		}
		return val, nil
	default:
		return reflect.New(nil), fmt.Errorf("cannot convert token to value: %v", token)
	}
}

func (token Token) SetValueOf(val reflect.Value) error {
	if !val.CanSet() {
		return nil
	}
	switch val.Kind() {
	case reflect.String:
		val.SetString(token.ToString())
		return nil
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int8, reflect.Int16:
		n, err := Ttoi(token)
		if err != nil {
			return err
		}
		val.SetInt(n)
		return err
	case reflect.Float32, reflect.Float64:
		f, err := Ttof(token)
		if err != nil {
			return err
		}
		val.SetFloat(f)
		return err
	case reflect.Bool:
		// TODO : this doesn't actually work it will parse any value starting with 'f' and/or 't' as boolean
		// so "boolean": faulty -> false
		switch token.Value[0] {
		case 't':
			val.SetBool(true)
		case 'f':
			val.SetBool(false)
		default:
			return Error(ErrInvalidValue, token.String())
		}
		return nil
	}

	if token.Type == Null {
		return nil
	}

	return fmt.Errorf("cannot convert token to value: %v", token)
}

func (token Token) ToValue() (reflect.Value, error) {
	switch token.Type {
	case String:
		val := reflect.New(ReflectTypeString).Elem()
		val.SetString(string(token.Value))
		return val, nil
	case Integer:
		val := reflect.New(ReflectTypeInteger).Elem()
		n, err := strconv.ParseInt(string(token.Value), 10, 64)
		if err != nil {
			return val, err
		}
		val.SetInt(n)
		return val, err
	case Float:
		val := reflect.New(ReflectTypeFloat).Elem()
		f, err := strconv.ParseFloat(string(token.Value), 64)
		if err != nil {
			return val, err
		}
		val.SetFloat(f)
		return val, err
	case Boolean:
		// TODO : this doesn't actually work it will parse any value starting with 'f' and/or 't' as boolean
		// so "boolean": faulty -> false
		val := reflect.New(ReflectTypeBool).Elem()
		switch token.Value[0] {
		case 't':
			val.SetBool(true)
		case 'f':
			val.SetBool(false)
		default:
			return val, Error(ErrInvalidValue, token.String())
		}
		return val, nil
	default:
		return reflect.New(nil), fmt.Errorf("cannot convert token to value: %v", token)
	}
}

func (token Token) IsEnding() bool {
	return token.Type.IsEnding()
}

func (token Token) String() string {
	return fmt.Sprintf("[%s](%s)", token.Value, token.Type)
}

func (token Token) ToString() string {
	return string(token.Value)
}

type Tokens []Token

func (tokens Tokens) Join(sep string) string {
	var buf bytes.Buffer
	for i, token := range tokens {
		buf.Write(token.Value)
		if i < len(tokens)-1 {
			buf.WriteString(sep)
		}
	}
	return buf.String()
}
