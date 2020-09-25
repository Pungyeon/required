package json

import "bytes"

var (
	Space     byte = ' '
	Tab       byte = '\t'
	NewLine   byte = '\n'
	Quotation byte = '"'
)

type TokenType string

func (t TokenType) IsEnding() bool {
	return t == ClosingBraceToken || t == ClosingCurlyToken ||
		t == ClosingBracketToken
}

const (
	UnknownToken        TokenType = "UNKNOWN"
	StringToken         TokenType = "STRING"
	KeyToken            TokenType = "KEY_TOKEN"
	ColonToken          TokenType = ":"
	CommaToken          TokenType = ","
	WhiteSpaceToken     TokenType = "WHITESPACE"
	OpenBraceToken      TokenType = "["
	ClosingBraceToken   TokenType = "]"
	OpenBracketToken    TokenType = "("
	ClosingBracketToken TokenType = ")"
	OpenCurlyToken      TokenType = "{"
	ClosingCurlyToken   TokenType = "}"
	FullStopToken       TokenType = "."
)

var TokenTypes = map[string]TokenType{
	"UNKNOWN":    UnknownToken,
	"STRING":     StringToken,
	"KEY_TOKEN":  KeyToken,
	":":          ColonToken,
	",":          CommaToken,
	"WHITESPACE": WhiteSpaceToken,
	"[":          OpenBraceToken,
	"]":          ClosingBraceToken,
	"(":          OpenBracketToken,
	")":          ClosingBracketToken,
	"{":          OpenCurlyToken,
	"}":          ClosingCurlyToken,
	".":          FullStopToken,
}

type Token struct {
	Value string
	Type  TokenType
}

func NewToken(b byte) Token {
	t, ok := TokenTypes[string(b)]
	if !ok {
		return Token{
			Value: string(b),
			Type:  UnknownToken,
		}
	}
	return Token{
		Value: string(b),
		Type:  t,
	}
}

func (token Token) IsEnding() bool {
	return token.Type.IsEnding()
}

func (token Token) String() string {
	return token.Value + ": " + string(token.Type)
}

type Tokens []Token

func (tokens Tokens) Join(sep string) string {
	var buf bytes.Buffer
	for i, token := range tokens {
		buf.WriteString(token.Value)
		if i < len(tokens)-1 {
			buf.WriteString(sep)
		}
	}
	return buf.String()
}
