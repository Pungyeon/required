package json

import (
	"errors"

	"github.com/Pungyeon/json-validation/pkg/token"
)

var (
	errInvalidJSONString = errors.New("invalid JSON string")
)

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
	index  int
	input  string
	output []token.Token
}

func (l *lexer) next() bool {
	l.index++
	return l.index < len(l.input)
}

func (l *lexer) value() byte {
	return l.input[l.index]
}

func (l *lexer) readNumber() (token.Token, error) {
	tokenType := token.Integer
	buf := []byte{l.value()}
	for l.next() {
		switch l.value() {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			buf = append(buf, l.value())
		case '.':
			tokenType = token.Float
			buf = append(buf, l.value())
		default:
			return token.Token{
				Value: string(buf),
				Type:  tokenType,
			}, nil
		}
	}
	return token.Token{
		Value: string(buf),
		Type:  tokenType,
	}, nil
}

func (l *lexer) readString() (token.Token, error) {
	var buf []byte
	for l.next() {
		if l.value() == token.Quotation {
			return token.Token{
				Value: string(buf),
				Type:  token.String,
			}, nil
		}
		buf = append(buf, l.value())
	}
	return token.Token{}, errors.New("unexpected end of file, trying to read string")
}
