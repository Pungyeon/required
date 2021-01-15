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

func BenchmarkLexerPerformance(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tokens, err := Lex(`{"foo": [1, 2, {"bar": 2}, true]}`)
		if err != nil {
			b.Fatal(err)
		}
		_ = tokens
	}
}

func BenchmarkLexerStreamPerformance(b *testing.B) {
	for i := 0; i < b.N; i++ {
		l := NewLexer([]byte(`{"foo": [1, 2, {"bar": 2}, true]}`))
		for l.Next() {
			_ = l.Current()
		}
	}
}
