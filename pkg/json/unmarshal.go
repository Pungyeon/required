package json

import (
	"io"

	"github.com/Pungyeon/required/pkg/lexer"
)

func Unmarshal(data []byte, v interface{}) error {
	return Parse(lexer.NewLexer(data), v)
}

type Decoder struct {
	r io.Reader
}

func NewDecoder(w io.Reader) *Decoder {
	return &Decoder{r: w}
}

func (d *Decoder) Decode(v interface{}) error {
	l, err := lexer.NewLexerReader(d.r)
	if err != nil {
		return err
	}
	return Parse(l, v)
}
