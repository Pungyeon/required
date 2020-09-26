package json

import (
	"encoding/json"
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
	Ding           int64
	Dong           string
	Float          float64
	Object         TestObject   `json:"object"`
	Array          []int        `json:"array"`
	StringSlice    []string     `json:"string_slice"`
	MultiDimension [][]int      `json:"multi_dimension"`
	ObjectArray    []TestObject `json:"obj_array"`
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
		"multi_dimension": [
			[1, 2, 3],
			[4, 5, 6]
		],
		"obj_array": [
			{"name": "lasse"},
			{"name": "basse"}
		]
	}`

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

	if ding.MultiDimension[1][0] != 4 {
		t.Fatalf("mismatch: (%v) != (%v)", ding.MultiDimension, [][]int{
			{1, 2, 3},
			{4, 5, 6}})
	}

	if len(ding.ObjectArray) != 2 {
		t.Fatalf("mismatch: (%d) != (%d)", len(ding.ObjectArray), 2)
	}

	if ding.ObjectArray[1].Name != "basse" {
		t.Fatalf("mismatch: (%v) != (%v)", ding.ObjectArray, []TestObject{
			{Name: "lasse"},
			{Name: "basse"},
		})
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

func TestParseMultiArray(t *testing.T) {
	tokens := Lex(`[
	[1, 2, 3],
	[4, 5, 6]
]`)
	if tokens.Join(";") != "[;[;1;,;2;,;3;];,;[;4;,;5;,;6;];]" {
		t.Fatal("oh no", tokens.Join(";"))
	}

	var obj [][]int
	if err := Parse(tokens, &obj); err != nil {
		t.Fatal(err)
	}

	if len(obj) != 2 {
		t.Fatal("length of object:", len(obj))
	}

	if len(obj[0]) != 3 {
		t.Fatal("length of object inner:", len(obj[0]))
	}

	if obj[1][0] != 4 {
		t.Fatal("omg it's not 4:", obj[1][0])
	}
}

func TestParseMultiStringArray(t *testing.T) {
	tokens := Lex(`[
	["1", "2", "3"],
	["4", "5", "6"]
]`)
	if tokens.Join(";") != "[;[;1;,;2;,;3;];,;[;4;,;5;,;6;];]" {
		t.Fatal("oh no", tokens.Join(";"))
	}

	var obj [][]string
	if err := Parse(tokens, &obj); err != nil {
		t.Fatal(err)
	}

	if len(obj) != 2 {
		t.Fatal("length of object:", len(obj))
	}

	if len(obj[0]) != 3 {
		t.Fatal("length of object inner:", len(obj[0]))
	}

	if obj[1][0] != "4" {
		t.Fatal("omg it's not 4:", obj[1][0])
	}
}

func TestParseObjectArray(t *testing.T) {
	tokens := Lex(`[
	{
		"name": "lasse"
	},
	{
		"name": "basse"
	}
]`)
	if tokens.Join(";") != "[;{;name;:;lasse;};,;{;name;:;basse;};]" {
		t.Fatal("oh no", tokens.Join(";"))
	}

	var obj []TestObject
	if err := Parse(tokens, &obj); err != nil {
		t.Fatal(err)
	}

	if len(obj) != 2 {
		t.Fatal("length of object:", len(obj))
	}

	if obj[1].Name != "basse" {
		t.Fatal("omg it's not basse:", obj[1].Name)
	}
}

func TestMapStringIntUnmarshal(t *testing.T) {
	tokens := Lex(`{"number": 1, "lumber": 13}`)
	var m map[string]int
	if err := Parse(tokens, &m); err != nil {
		t.Fatal(err)
	}
	if m["number"] != 1 {
		t.Fatal("map parsed incorrectly:", m)
	}

	if m["lumber"] != 13 {
		t.Fatal("map parsed incorrectly:", m)
	}
}

func TestMapStringStringUnmarshal(t *testing.T) {
	tokens := Lex(`{"number": "1", "lumber": "13"}`)
	var m map[string]string
	if err := Parse(tokens, &m); err != nil {
		t.Fatal(err)
	}
	if m["number"] != "1" {
		t.Fatal("map parsed incorrectly:", m)
	}

	if m["lumber"] != "13" {
		t.Fatal("map parsed incorrectly:", m)
	}
}

func BenchmarkStdUnmarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var ding Ding
		if err := json.Unmarshal([]byte(sample), &ding); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPkgUnmarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var ding Ding
		if err := Unmarshal([]byte(sample), &ding); err != nil {
			b.Fatal(err)
		}
	}
}
