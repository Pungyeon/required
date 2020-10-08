package json

import (
	"encoding/json"
	"fmt"
	"testing"
)

type TestObject struct {
	Name string `json:"name"`
}

type Ding struct {
	Ding           int64
	Boolean        bool
	Dong           string
	Float          float64
	Object         *TestObject  `json:"object,required"`
	Array          []int64      `json:"array"`
	StringSlice    []string     `json:"string_slice"`
	MultiDimension [][]int      `json:"multi_dimension"`
	ObjectArray    []TestObject `json:"obj_array"`
	MapObject      map[string]int
}

var sample = `{
		"ding": 1,
		"boolean": true,
		"dong": "hello",
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
		],
		"map_object": {
			"number": 1,
			"lumber": 13
		},
		"float": 3.2
	}`

func TestLexer(t *testing.T) {
	tokens := Lex(`{"foo": [1, 2, {"bar": 2}, true]}`)

	result := tokens.Join(";")
	expected := "{;foo;:;[;1;,;2;,;{;bar;:;2;};,;true;];}"

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

func TestParsePrimitive(t *testing.T) {
	var v int64
	if err := Parse(Lex(`1`), &v); err != nil {
		t.Fatal(err)
	}
	if v != 1 {
		t.Fatal("v not equal 1:", v)
	}
}

func TestParseArrayInStruct(t *testing.T) {
	type Thing struct {
		Array []int64
	}
	tokens := Lex(`{"array": [1, 2, 3, 4]}`)

	var obj Thing
	if err := Parse(tokens, &obj); err != nil {
		t.Fatal(err)
	}

	if len(obj.Array) != 4 {
		t.Fatal(len(obj.Array))
	}

	if obj.Array[2] != 3 {
		t.Fatal("expected 3:", obj.Array[2])
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

	if len(obj) != 4 {
		t.Fatal(len(obj))
	}

	if obj[2] != 3 {
		t.Fatal("expected 3:", obj[2])
	}
}

func TestParseFloatArray(t *testing.T) {
	tokens := Lex("[1.1, 2.2, 3.3, 4.4]")
	if tokens.Join(";") != "[;1.1;,;2.2;,;3.3;,;4.4;]" {
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

func testParse(t *testing.T, tokens Tokens, v interface{}) {
	if err := Parse(tokens, v); err != nil {
		t.Fatal(err)
	}
}

func TestParseAsReflectValue(t *testing.T) {

	tt := []struct {
		name  string
		check func() bool
	}{
		{name: "string", check: func() bool {
			var v string
			testParse(t, Lex(`"lasse"`), &v)
			return v == "lasse"
		}},
		{name: "int64", check: func() bool {
			var v int64
			testParse(t, Lex(`234`), &v)
			return v == 234
		}},
		{name: "int32", check: func() bool {
			var v int32
			testParse(t, Lex(`234`), &v)
			return v == 234
		}},
		{name: "int", check: func() bool {
			var v int
			testParse(t, Lex(`234`), &v)
			return v == 234
		}},
		{name: "float64", check: func() bool {
			var v float64
			testParse(t, Lex(`42.2`), &v)
			return v == 42.2
		}},
		{name: "float32", check: func() bool {
			var v float32
			testParse(t, Lex(`42.2`), &v)
			return v == 42.2
		}},
		{name: "test_object", check: func() bool {
			var v TestObject
			testParse(t, Lex(`{"name": "lasse"}`), &v)
			return v.Name == "lasse"
		}},
		{name: "array", check: func() bool {
			var v []string
			testParse(t, Lex(`["name", "lasse"]`), &v)
			return v[1] == "lasse"
		}},
		{name: "array", check: func() bool {
			var v interface{}
			testParse(t, Lex(`"lasse"`), &v)
			return v.(string) == "lasse"
		}},
		{name: "array", check: func() bool {
			var v interface{}
			testParse(t, Lex(`{"name": "lasse"}`), &v)
			return v.(map[string]interface{})["name"] == "lasse"
		}},
		{name: "interface_array", check: func() bool {
			var v []interface{}
			testParse(t, Lex(`["name", "lasse"]`), &v)
			return v != nil &&
				v[0].(string) == "name"
		}},
		{name: "ding_object", check: func() bool {
			var ding Ding
			testParse(t, Lex(sample), &ding)
			return ding.Ding == 1 &&
				ding.Dong == "hello" &&
				ding.Boolean == true &&
				ding.Object.Name == "lasse" &&
				ding.Array[2] == 3 &&
				ding.StringSlice[2] == "3" &&
				ding.MultiDimension[1][2] == 6 &&
				ding.ObjectArray[1].Name == "basse" &&
				ding.MapObject["lumber"] == 13 &&
				ding.Float == 3.2
		}},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.check() {
				t.Fatalf("check failed on test: %s", tc.name)
			}
		})
	}
}

func TestParsePointer(t *testing.T) {
	tokens := Lex(`{
		"object": {
			"name": "lasse"
		},
	}`)

	var ding Ding
	if err := Parse(tokens, &ding); err != nil {
		t.Fatal(err)
	}

	fmt.Println(ding)
	if ding.Object.Name != "lasse" {
		t.Fatal("oh no")
	}
}

func TestParseInterfaceString(t *testing.T) {
	tokens := Lex(`"lasse"`)

	var ding interface{}
	if err := Parse(tokens, &ding); err != nil {
		t.Fatal(err)
	}

	fmt.Println(ding)
	if ding.(string) != "lasse" {
		t.Fatal("oh no")
	}
}

func TestMapFollowedBy(t *testing.T) {
	tokens := Lex(`{
	"map_object": {
		"number": 1,
			"lumber": 13
	},
	"float": 3.2
}`)
	var ding Ding
	if err := Parse(tokens, &ding); err != nil {
		t.Fatal(err)
	}
	if ding.MapObject["number"] != 1 ||
		ding.MapObject["lumber"] != 13 ||
		ding.Float != 3.2 {
		t.Fatal("Unexpected result:", ding)
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
