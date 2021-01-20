package token

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
)

var (
	Space     byte = ' '
	Tab       byte = '\t'
	NewLine   byte = '\n'
	Quotation byte = '"'
)

var (
	ReflectTypeString    = reflect.TypeOf("")
	ReflectTypeInteger   = reflect.TypeOf(1)
	ReflectTypeFloat     = reflect.TypeOf(3.2)
	ReflectTypeInterface = reflect.ValueOf(map[string]interface{}{}).Type().Elem()
	ReflectTypeBool      = reflect.TypeOf(true)
)

type TokenType int

func (t TokenType) IsEnding() bool {
	return t == ClosingBrace || t == ClosingCurly ||
		t == ClosingBracket
}

const (
	Unknown TokenType = iota
	Integer
	Float
	String
	Null
	Key
	Colon
	Comma
	WhiteSpace
	OpenBrace
	ClosingBrace
	OpenBracket
	ClosingBracket
	OpenCurly
	ClosingCurly
	FullStop
	Boolean
)

var TokenTypes = map[byte]TokenType{
	// "UNKNOWN":    Unknown,
	// "BOOLEAN":    Boolean,
	// "INTEGER":    Integer,
	// "FLOAT":      Float,
	// "STRING":     String,
	// "NULL":       Null,
	// "KEY_TOKEN":  Key,
	// "WHITESPACE": WhiteSpace,
	':': Colon,
	',': Comma,
	'[': OpenBrace,
	']': ClosingBrace,
	'(': OpenBracket,
	')': ClosingBracket,
	'{': OpenCurly,
	'}': ClosingCurly,
	'.': FullStop,
}

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
	t, ok := TokenTypes[b[i]]
	if !ok {

		return Token{
			Value: b[i : i+1], // should we even allocate here?
			Type:  Unknown,
		}
	}
	return Token{
		Value: b[i : i+1],
		Type:  t,
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
		n, err := strconv.ParseInt(string(token.Value), 10, 64)
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
		if token.Value[0] == 't' {
			val.SetBool(true)
		} else {
			val.SetBool(false)
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
	switch token.Type {
	case Null:
		return nil // don't set anything
	case String:
		val.SetString(token.ToString())
		return nil
	case Integer:
		n, err := strconv.ParseInt(token.ToString(), 10, 64)
		if err != nil {
			return err
		}
		val.SetInt(n)
		return err
	case Float:
		f, err := strconv.ParseFloat(token.ToString(), 64)
		if err != nil {
			return err
		}
		val.SetFloat(f)
		return err
	case Boolean:
		if token.Value[0] == 't' {
			val.SetBool(true)
		} else {
			val.SetBool(false)
		}
		return nil
	default:
		return fmt.Errorf("cannot convert token to value: %v", token)
	}
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
		val := reflect.New(ReflectTypeBool).Elem()
		if token.Value[0] == 't' {
			val.SetBool(true)
		} else {
			val.SetBool(false)
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
	return fmt.Sprintf("%s: %d", token.Value, token.Type)
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
