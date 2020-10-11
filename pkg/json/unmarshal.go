package json

func Unmarshal(data []byte, v interface{}) error {
	tokens, err := Lex(string(data))
	if err != nil {
		return err
	}
	return Parse(tokens, v)
}
