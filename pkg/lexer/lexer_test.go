package lexer

import (
	"testing"
)

func TestLexer(t *testing.T) {
	tokens, err := Lex(`{"foo": [1, 2, {"bar": 2}, true]}`)
	if err != nil {
		t.Fatal(err)
	}

	result := tokens.Join(";")
	expected := "{;foo;:;[;1;,;2;,;{;bar;:;2;};,;true;];}"

	if result != expected {
		t.Fatalf("%v != %v", result, expected)
	}
}
