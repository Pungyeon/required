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

type Lexer interface {
	EOF() bool
	Previous() token.Token
	Current() token.Token
	Next() bool
}

func NewLexer(input string) Lexer {
	return &lexer{
		input: input,
		index: -1,
	}
}

func (l *lexer) EOF() bool {
	return l.index >= len(l.input)
}

func (l *lexer) Previous() token.Token {
	return l.previous
}

func (l *lexer) Current() token.Token {
	return l.current
}

func (l *lexer) Next() bool {
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
		l.assign(token.Token{Value: "true", Type: token.Boolean})
		return true
	case 'f':
		l.index += len("alse")
		return l.assign(token.Token{Value: "false", Type: token.Boolean})
	case 'n':
		l.index += len("ull")
		return l.assign(token.Token{Value: "null", Type: token.Null})
	default:
		return l.assign(token.NewToken(l.value()))
	}
}

func Lex(input string) (token.Tokens, error) {
	l := &lexer{
		input: input,
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
			l.output = append(l.output, token.Token{Value: "true", Type: token.Boolean})
			l.index += len("rue")
		case 'f':
			l.output = append(l.output, token.Token{Value: "false", Type: token.Boolean})
			l.index += len("alse")
		case 'n':
			l.output = append(l.output, token.Token{Value: "null", Type: token.Null})
			l.index += len("ull")
		default:
			l.output = append(l.output, token.NewToken(l.value()))
		}
	}
	return l.output, nil
}

type lexer struct {
	current  token.Token
	previous token.Token
	index    int
	input    string
	output   []token.Token
}

func (l *lexer) next() bool {
	l.index++
	return l.index < len(l.input)
}

func (l *lexer) value() byte {
	return l.input[l.index]
}

func (l *lexer) assign(t token.Token) bool {
	l.previous = l.current
	l.current = t
	return true
}

func (l *lexer) readNumber() (token.Token, error) {
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

func (l *lexer) readString() (token.Token, error) {
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
