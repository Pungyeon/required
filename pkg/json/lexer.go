package json

import "errors"

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
