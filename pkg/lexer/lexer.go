package lexer

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/Pungyeon/required/pkg/token"
)

type Lexer struct {
	current  token.Token
	previous token.Token
	index    int
	input    []byte
	err      error
	stack *Stack
}

func NewLexerReader(r io.Reader) (*Lexer, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return NewLexer(data), nil
}

func NewLexer(input []byte) *Lexer {
	return &Lexer{
		input: input,
		index: -1,
		stack: NewStack(10),
	}
}

func (l *Lexer) GetError() error {
	if l.stack.index != 0 {
		return token.Error(token.ErrMissingBrace, fmt.Sprintf("expected '%s', but not found: %s", string(l.stack.Pop()), l.input[:]))
	}
	return l.err
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

func (l *Lexer) CurrentLine() string {
	i, j := l.index-1, l.index
	if i < 0 {
		i = 0
	}
	for i > 0 {
		if l.input[i] == '\n' {
			break
		}
		i--
	}
	for j < len(l.input) {
		if l.input[i] == '\n' {
			return string(l.input[i:j])
		}
		j++
	}
	return string(l.input[i:])
}

func (l *Lexer) SkipValue() []byte {
	if l.index == -1 {
		l.index = 0
	}
	l.skipWhile(':')
	l.skipWhitespace()
	var (
		stack   int
		start   = l.index
		opening byte
	)
	if l.value() == '"' {
		t, err := l.readString()
		if err != nil {
			return nil
		}
		return t.Value
	}

	closing, ok := token.BraceOpposites[l.value()]
	if !ok {
		closing = ','
	} else {
		opening = l.value()
	}

	for l.next() {
		if l.value() == opening {
			stack++
		}
		if l.value() == closing {
			if stack == 0 {
				return l.input[start : l.index+1]
			} else {
				stack--
			}
		}
	}
	if l.input[start] == '"' {

	}

	if l.input[l.index-1] == '}' {
		return l.input[start : l.index-1]
	}
	if l.index == len(l.input) {
		return l.input[start:l.index]
	}
	return l.input[start : l.index+1]
}

func (l *Lexer) skipWhitespace() {
	l.skipWhile(' ')
}

func (l *Lexer) skipWhile(b byte) {
	for !l.EOF() && l.value() == b {
		if !l.next() {
			return
		}
	}
}

func (l *Lexer) skipTo(b byte) {
	for l.next() {
		if l.value() == b {
			return
		}
	}
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
			l.err = token.Error(token.ErrInvalidJSON, string(l.input[:l.index]))
			return false
		}
		return l.assign(t)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		t, err := l.readNumber()
		if err != nil {
			l.err = token.Error(token.ErrInvalidJSON, string(l.input[:l.index]))
			return false
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

func (l *Lexer) next() bool {
	l.index++
	return l.index < len(l.input)
}

func (l *Lexer) value() byte {
	return l.input[l.index]
}

func (l *Lexer) assign(t token.Token) bool {
	if t.Type.IsOpening() {
		l.stack.Push(t.Value[0])
	}
	if t.Type.IsEnding() {
		if token.BraceOpposites[l.stack.Pop()] != t.Value[0] {
			l.err = token.Error(token.ErrUnmatchedBrace, string(l.input[:l.index]))
			return false
		}
	}
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
