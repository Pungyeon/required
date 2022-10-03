package lexer

import (
	"io"
	"io/ioutil"

	"github.com/Pungyeon/required/pkg/token"
)

type Lexer struct {
	index int
	input []byte
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

func (l *Lexer) Previous() string {
	return string(l.input[max(0, l.index-100):l.index])
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (l *Lexer) EOF() bool {
	return l.index >= len(l.input)
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
		// TODO what in the world???
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

func (l *Lexer) Next() (token.Token, error) {
	if !l.next() {
		return token.Empty, l.isValid()
	}
	switch l.value() {
	case token.Space, token.Tab, token.NewLine:
		return l.Next()
	case token.Quotation:
		return l.readString()
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		t := l.readNumber()
		l.index--
		return t, nil
	case 't':
		l.index += len("rue")
		return token.Token{Value: TRUE, Type: token.Boolean}, nil
	case 'f':
		l.index += len("alse")
		return token.Token{Value: FALSE, Type: token.Boolean}, nil
	case 'n':
		l.index += len("ull")
		return token.Token{Value: NULL, Type: token.Null}, nil
	default:
		t := token.NewToken(l.input, l.index)
		if t.Type.IsOpening() {
			l.stack.Push(l.value())
		}
		if t.Type.IsEnding() {
			opposite := l.stack.Pop()
			if token.BraceOpposites[opposite] != t.Value[0] {
				return t, token.Error(token.ErrUnmatchedBrace, string(l.input[:l.index]))
			}
		}
		return t, nil
	}
}

func (l *Lexer) next() bool {
	l.index++
	return l.index < len(l.input)
}

func (l *Lexer) value() byte {
	return l.input[l.index]
}

func (l *Lexer) previous() byte {
	return l.input[l.index-1]
}

func (l *Lexer) readNumber() token.Token {
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
			}
		}
	}
	return token.Token{
		Value: l.input[start:l.index],
		Type:  tokenType,
	}
}

func (l *Lexer) readString() (token.Token, error) {
	start := l.index + 1
	for l.next() {
		if l.value() == token.Quotation && l.previous() != token.Escape {
			return token.Token{
				Value: l.input[start:l.index],
				Type:  token.String,
			}, nil
		}
	}
	return token.Empty, token.Error(token.ErrInvalidJSON, string(l.input[:l.index]))
}

func (l *Lexer) isValid() error {
	if l.stack.IsEmpty() {
		return io.EOF
	}
	return token.Error(token.ErrUnmatchedBrace, l.Previous())
}
