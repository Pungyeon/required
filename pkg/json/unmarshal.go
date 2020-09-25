package json

func Unmarshal(data []byte, v interface{}) error {
	return Parse(Lex(string(data)), v)
}
