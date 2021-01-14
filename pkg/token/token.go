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

type TokenType string

func (t TokenType) IsEnding() bool {
	return t == ClosingBrace || t == ClosingCurly ||
		t == ClosingBracket
}

const (
	Unknown        TokenType = "UNKNOWN"
	Integer        TokenType = "INTEGER"
	Float          TokenType = "FLOAT"
	String         TokenType = "STRING"
	Null           TokenType = "NULL"
	Key            TokenType = "KEY_TOKEN"
	Colon          TokenType = ":"
	Comma          TokenType = ","
	WhiteSpace     TokenType = "WHITESPACE"
	OpenBrace      TokenType = "["
	ClosingBrace   TokenType = "]"
	OpenBracket    TokenType = "("
	ClosingBracket TokenType = ")"
	OpenCurly      TokenType = "{"
	ClosingCurly   TokenType = "}"
	FullStop       TokenType = "."
	Boolean        TokenType = "BOOLEAN"
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

type Token struct {
	Value interface{}
	Type  TokenType
}

func NewToken(b byte) Token {
	t, ok := TokenTypes[b]
	if !ok {

		return Token{
			Value: b,
			Type:  Unknown,
		}
	}
	return Token{
		Value: b,
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
		val.SetString(token.Value.(string))
		return val, nil
	case Integer:
		n, err := strconv.ParseInt(token.Value.(string), 10, 64)
		if err != nil {
			return val, err
		}
		val.SetInt(n)
		return val, err
	case Float:
		f, err := strconv.ParseFloat(token.Value.(string), 64)
		if err != nil {
			return val, err
		}
		val.SetFloat(f)
		return val, err
	case Boolean:
		if token.Value.(string)[0] == 't' {
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
	switch token.Type {
	case String:
		val.SetString(token.Value.(string))
		return nil
	case Integer:
		n, err := strconv.ParseInt(token.Value.(string), 10, 64)
		if err != nil {
			return err
		}
		val.SetInt(n)
		return err
	case Float:
		f, err := strconv.ParseFloat(token.Value.(string), 64)
		if err != nil {
			return err
		}
		val.SetFloat(f)
		return err
	case Boolean:
		if token.Value.(string)[0] == 't' {
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
		val.SetString(token.Value.(string))
		return val, nil
	case Integer:
		val := reflect.New(ReflectTypeInteger).Elem()
		n, err := strconv.ParseInt(token.Value.(string), 10, 64)
		if err != nil {
			return val, err
		}
		val.SetInt(n)
		return val, err
	case Float:
		val := reflect.New(ReflectTypeFloat).Elem()
		f, err := strconv.ParseFloat(token.Value.(string), 64)
		if err != nil {
			return val, err
		}
		val.SetFloat(f)
		return val, err
	case Boolean:
		val := reflect.New(ReflectTypeBool).Elem()
		if token.Value.(string)[0] == 't' {
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
	return token.Value.(string) + ": " + string(token.Type)
}

type Tokens []Token

func (tokens Tokens) Join(sep string) string {
	var buf bytes.Buffer
	for i, token := range tokens {
		switch token.Value.(type) {
		case string:
			buf.WriteString(token.Value.(string))
			if i < len(tokens)-1 {
				buf.WriteString(sep)
			}
		case byte:
			buf.WriteString(string(token.Value.(byte)))
			if i < len(tokens)-1 {
				buf.WriteString(sep)
			}
		}
	}
	return buf.String()
}
