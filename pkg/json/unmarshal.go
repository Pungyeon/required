package json

import "github.com/Pungyeon/required/pkg/lexer"

func Unmarshal(data []byte, v interface{}) error {
	//tokens, err := lexer.Lex(string(data))
	//if err != nil {
	//	return err
	//}
	return Parse(lexer.NewLexer(data), v)
}

func DecodeJSON(data []byte, v interface{}) error {
	return Decode(lexer.NewLexer(data), v)
}
