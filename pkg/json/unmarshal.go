package json

import "github.com/Pungyeon/json-validation/pkg/lexer"

func Unmarshal(data []byte, v interface{}) error {
	tokens, err := lexer.Lex(string(data))
	if err != nil {
		return err
	}
	return Parse(tokens, v)
}
