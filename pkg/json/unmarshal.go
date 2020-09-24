package json

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
)

type TokenType string

const (
	UnknownToken        TokenType = "UNKNOWN"
	StringToken         TokenType = "STRING"
	KeyToken            TokenType = "KEY_TOKEN"
	ColonToken          TokenType = ":"
	CommaToken          TokenType = ","
	WhiteSpaceToken     TokenType = "WHITESPACE"
	OpenBraceToken      TokenType = "["
	ClosingBraceToken   TokenType = "]"
	OpenBracketToken    TokenType = "("
	ClosingBracketToken TokenType = ")"
	OpenCurlyToken      TokenType = "{"
	ClosingCurlyToken   TokenType = "}"
)

var TokenTypes = map[string]TokenType{
	"UNKNOWN":    UnknownToken,
	"STRING":     StringToken,
	"KEY_TOKEN":  KeyToken,
	":":          ColonToken,
	",":          CommaToken,
	"WHITESPACE": WhiteSpaceToken,
	"[":          OpenBraceToken,
	"]":          ClosingBraceToken,
	"(":          OpenBracketToken,
	")":          ClosingBracketToken,
	"{":          OpenCurlyToken,
	"}":          ClosingCurlyToken,
}

type Token struct {
	Value string
	Type  TokenType
}

func NewToken(b byte) Token {
	t, ok := TokenTypes[string(b)]
	if !ok {
		return Token{
			Value: string(b),
			Type:  UnknownToken,
		}
	}
	return Token{
		Value: string(b),
		Type:  t,
	}
}

func (token Token) String() string {
	return token.Value + ": " + string(token.Type)
}

type Tokens []Token

func (tokens Tokens) Join(sep string) string {
	var buf bytes.Buffer
	for i, token := range tokens {
		buf.WriteString(token.Value)
		if i < len(tokens)-1 {
			buf.WriteString(sep)
		}
	}
	return buf.String()
}

func Lex(input string) Tokens {
	l := &lexer{
		input: input,
		index: -1,
	}

	for l.next() {
		switch l.value() {
		case Space, Tab, NewLine:
			continue
		case Quotation:
			str, err := l.readString()
			if err != nil {
				panic(err)
			}
			l.output = append(l.output, Token{
				Value: str,
				Type:  StringToken,
			})
		default:
			l.output = append(l.output, NewToken(l.value()))
		}
	}
	return l.output
}

type lexer struct {
	index  int
	input  string
	output []Token
}

func (l *lexer) next() bool {
	l.index++
	return l.index < len(l.input)
}

func (l *lexer) value() byte {
	return l.input[l.index]
}

func (l *lexer) readString() (string, error) {
	//l.next() // skip current quotation
	var buf []byte
	for l.next() {
		if l.value() == Quotation {
			return string(buf), nil
		}
		buf = append(buf, l.value())
	}
	return "", errors.New("unexpected end of file, trying to read string")
}

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
