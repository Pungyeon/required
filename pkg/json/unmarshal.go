package json

import "github.com/Pungyeon/required/pkg/lexer"

func Unmarshal(data []byte, v interface{}) error {
	return Parse(lexer.NewLexer(data), v)
}
