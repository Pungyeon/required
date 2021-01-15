package lexer

import (
	"errors"

	"github.com/Pungyeon/required/pkg/token"
)

var (
	errInvalidJSONString = errors.New("invalid JSON string")
)

type Result struct {
	Token token.Token
	Error error
}

type ILexer interface {
	EOF() bool
	Previous() token.Token
	Current() token.Token
	Next() bool
}

func NewLexer(input []byte) *Lexer {
	return &Lexer{
		input: input,
		index: -1,
	}
}

func (l *Lexer) EOF() bool {
	return l.index >= len(l.input)
}

func (l *Lexer) Previous() token.Token {
	return l.previous
}

func (l *Lexer) Current() token.Token {
	return l.current
}

var (
	TRUE  = []byte("true")
	FALSE = []byte("false")
	NULL  = []byte("null")
)

func (l *Lexer) Next() bool {
	if !l.next() {
		return false // should be eof error?
	}
	switch l.value() {
	case token.Space, token.Tab, token.NewLine:
		return l.Next()
	case token.Quotation:
		t, err := l.readString()
		if err != nil {
			panic(errInvalidJSONString) // TODO : no panic plox
		}
		return l.assign(t)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		t, err := l.readNumber()
		if err != nil {
			panic(errInvalidJSONString) // TODO : no panic plox
		}
		l.index--
		return l.assign(t)
	case 't':
		l.index += len("rue")
		l.assign(token.Token{Value: TRUE, Type: token.Boolean})
		return true
	case 'f':
		l.index += len("alse")
		return l.assign(token.Token{Value: FALSE, Type: token.Boolean})
	case 'n':
		l.index += len("ull")
		return l.assign(token.Token{Value: NULL, Type: token.Null})
	default:
		return l.assign(token.NewToken(l.input, l.index))
	}
}

func Lex(input string) (token.Tokens, error) {
	l := &Lexer{
		input: []byte(input),
		index: -1,
	}

	for l.next() {
		switch l.value() {
		case token.Space, token.Tab, token.NewLine:
			continue
		case token.Quotation:
			t, err := l.readString()
			if err != nil {
				return token.Tokens{}, errInvalidJSONString
			}
			l.output = append(l.output, t)
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			t, err := l.readNumber()
			if err != nil {
				return token.Tokens{}, errInvalidJSONString
			}
			l.output = append(l.output, t)
			l.index--
		case 't':
			l.output = append(l.output, token.Token{Value: TRUE, Type: token.Boolean})
			l.index += len("rue")
		case 'f':
			l.output = append(l.output, token.Token{Value: FALSE, Type: token.Boolean})
			l.index += len("alse")
		case 'n':
			l.output = append(l.output, token.Token{Value: NULL, Type: token.Null})
			l.index += len("ull")
		default:
			l.output = append(l.output, token.NewToken(l.input, l.index))
		}
	}
	return l.output, nil
}

type Lexer struct {
	current  token.Token
	previous token.Token
	index    int
	input    []byte
	output   []token.Token
}

func (l *Lexer) next() bool {
	l.index++
	return l.index < len(l.input)
}

func (l *Lexer) value() byte {
	return l.input[l.index]
}

func (l *Lexer) assign(t token.Token) bool {
	l.previous = l.current
	l.current = t
	return true
}

func (l *Lexer) readNumber() (token.Token, error) {
	tokenType := token.Integer
	start := l.index
	for l.next() {
		switch l.value() {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		case '.':
			tokenType = token.Float
		default:
			return token.Token{
				Value: l.input[start:l.index],
				Type:  tokenType,
			}, nil
		}
	}
	return token.Token{
		Value: l.input[start:l.index],
		Type:  tokenType,
	}, nil
}

func (l *Lexer) readString() (token.Token, error) {
	start := l.index + 1
	for l.next() {
		if l.value() == token.Quotation {
			return token.Token{
				Value: l.input[start:l.index],
				Type:  token.String,
			}, nil
		}
	}
	return token.Token{}, errors.New("unexpected end of file, trying to read string")
}
