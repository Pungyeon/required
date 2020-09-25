package json

import (
	"fmt"
	"testing"
)

type TheThing struct {
	Ding
}

type TestObject struct {
	Name string `json:"name"`
}

type Ding struct {
	Ding        int64      `json:"ding"`
	Dong        string     `json:"dong"`
	Float       float64    `json:"float"`
	Object      TestObject `json:"object"`
	Array       []int      `json:"array"`
	StringSlice []string   `json:"string_slice"`
}

var sample = `{
		"ding": 1,
		"dong": "hello",
		"float": 3.2,
		"object": {
			"name": "lasse"
		},
		"array": [1, 2, 3],
		"string_slice": ["1", "2", "3"],
	}`

//"multidimensional_array": [
//[1, 2, 3],
//[4, 5, 6]
//]
func TestLexer(t *testing.T) {
	tokens := Lex(`{"foo": [1, 2, {"bar": 2}]}`)

	result := tokens.Join(";")
	expected := "{;foo;:;[;1;,;2;,;{;bar;:;2;};];}"

	if result != expected {
		t.Fatalf("%v != %v", result, expected)
	}
}

func TestParserSimple(t *testing.T) {
	var obj TestObject
	if err := Parse(Lex(`{"name": "lasse"}`), &obj); err != nil {
		t.Fatal(err)
	}
	if obj.Name != "lasse" {
		t.Fatal("not lasse:", obj.Name)
	}
}

func TestParseComplex(t *testing.T) {
	var ding Ding
	if err := Parse(Lex(sample), &ding); err != nil {
		t.Fatal(err)
	}

	if ding.Ding != 1 {
		t.Fatalf("mismatch: (%d) != (%d)", ding.Ding, 1)
	}

	if ding.Dong != "hello" {
		t.Fatalf("mismatch: (%s) != (%s)", ding.Dong, "hello")
	}

	if ding.Float != 3.2 {
		t.Fatalf("mismatch: (%f) != (%f)", ding.Float, 3.2)
	}

	if ding.Object.Name != "lasse" {
		t.Fatalf("mismatch: (%s) != (%s)", ding.Object.Name, "lasse")
	}

	if len(ding.Array) != 3 {
		t.Fatalf("mismatch: (%d) != (%d)", len(ding.Array), 3)
	}

	if ding.Array[2] != 3 {
		t.Fatalf("mismatch: (%v) != (%v)", ding.Array, []int{1, 2, 3})
	}

	if len(ding.StringSlice) != 3 {
		t.Fatalf("mismatch: (%d) != (%d)", len(ding.StringSlice), 3)
	}

	if ding.StringSlice[2] != "3" {
		t.Fatalf("mismatch: (%v) != (%v)", ding.StringSlice, []string{"1", "2", "3"})
	}
}

func TestParseArray(t *testing.T) {
	tokens := Lex("[1, 2, 3, 4]")
	if tokens.Join(";") != "[;1;,;2;,;3;,;4;]" {
		t.Fatal("oh no", tokens.Join(";"))
	}

	var obj []int
	if err := Parse(tokens, &obj); err != nil {
		t.Fatal(err)
	}

	fmt.Println(obj)
	if len(obj) != 4 {
		t.Fatal(len(obj))
	}

	if obj[2] != 3 {
		t.Fatal("expected 3:", obj[2])
	}
}

func TestParseFloatArray(t *testing.T) {
	tokens := Lex("[1.1, 2.2, 3.3, 4.4]")
	if tokens.Join(";") != "[;1;.;1;,;2;.;2;,;3;.;3;,;4;.;4;]" {
		t.Fatal("oh no", tokens.Join(";"))
	}

	var obj []float64
	if err := Parse(tokens, &obj); err != nil {
		t.Fatal(err)
	}

	if len(obj) != 4 {
		t.Fatal(len(obj))
	}

	if obj[2] != 3.3 {
		t.Fatal("expected 3.3:", obj[2])
	}
}

//func BenchmarkStdUnmarshal(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		var ding Ding
//		if err := json.Unmarshal(sample, &ding); err != nil {
//			b.Fatal(err)
//		}
//	}
//}
//
//func BenchmarkPkgUnmarshal(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		var ding Ding
//		if err := Unmarshal(sample, &ding); err != nil {
//			b.Fatal(err)
//		}
//	}
//}
