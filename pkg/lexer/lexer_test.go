package lexer

import (
	"encoding/json"
	"fmt"
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

func TestScanValue(t *testing.T) {
	lexer := NewLexer([]byte(`{
	"foo": {
		"bar": "value" 
	},
	"float": 3.2
}`))
	lexer.skipTo(':')
	lexer.next()
	val := lexer.SkipValue()

	if string(val) != `{
		"bar": "value" 
	}` {
		t.Fatal(string(val))
	}
	lexer.skipTo(':')
	lexer.next()

	fmt.Println(string(lexer.input[lexer.index:]))
	fmt.Println("here we are")

	val = lexer.SkipValue()

	var floaty float64
	if err := json.Unmarshal(val, &floaty); err != nil {
		t.Fatal(err)
	}
	if floaty != 3.2 {
		t.Fatal(string(val))
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
