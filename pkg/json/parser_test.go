package json

import (
	"encoding/json"
	"fmt"
	"reflect"
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
	Object         *TestObject  `json:"object"`
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

func TestParseAsReflectValue(t *testing.T) {
	var val reflect.Value
	var i interface{}
	var Int64 int64
	var Int32 int32
	var Float64 float64
	var Float32 float32

	tt := []struct {
		name   string
		tokens Tokens
		Type   reflect.Type
		check  func() bool
	}{
		{"string", Lex(`"lasse"`), reflectTypeString, func() bool { return val.String() == "lasse" }},
		{"in64", Lex(`234`), reflect.TypeOf(Int64), func() bool { return val.Int() == 234 }},
		{"in32", Lex(`234`), reflect.TypeOf(Int32), func() bool { return val.Int() == 234 }},
		{"int", Lex(`13`), reflectTypeInteger, func() bool { return val.Int() == 13 }},
		{"float", Lex(`42.2`), reflectTypeFloat, func() bool { return val.Float() == 42.2 }},
		{"float64", Lex(`42.2`), reflect.TypeOf(Float64), func() bool { return val.Float() == 42.2 }},
		{"float32", Lex(`42.2`), reflect.TypeOf(Float32), func() bool { return val.Interface().(float32) == 42.2 }},
		{"test_object", Lex(`{"name": "lasse"}`),
			reflect.TypeOf(TestObject{}),
			func() bool { return val.Interface().(TestObject).Name == "lasse" },
		},
		{"array", Lex(`["name", "lasse"]`),
			reflect.TypeOf([]string{}),
			func() bool { return val.Interface().([]string)[1] == "lasse" },
		},
		{"interface_string", Lex(`"lasse"`),
			reflect.TypeOf(i),
			func() bool { return val.Interface().(string) == "lasse" },
		},
		{"interface_object", Lex(`{"name": "lasse"}`),
			reflect.TypeOf(i),
			func() bool {
				return val.Interface().(map[string]interface{})["name"].(string) == "lasse"
			},
		},
		{"interface_array", Lex(`["name", "lasse"]`),
			reflect.TypeOf(i),
			func() bool { return val.Interface().([]interface{}) != nil },
		},
		{"ding_object", Lex(sample),
			reflect.TypeOf(Ding{}),
			func() bool {
				ding := val.Interface().(Ding)
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
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			//p := &parser{
			//	tokens: tc.tokens,
			//	index:  -1,
			//}
			//
			//var err error
			//val, err = p.parse(tc.Type)
			//if err != nil {
			//	t.Error(tc.tokens)
			//	t.Fatal(err)
			//}
			if !tc.check() {
				t.Fatalf("%#v", val)
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

func TestPointers(t *testing.T) {
	var to TestObject
	fmt.Println("name:", to.Name)

	v := reflect.ValueOf(&to).Elem()
	field := v.FieldByName("Name")
	field.SetString("Lasse")
	fmt.Println("name:", to.Name)
	t.Error()

	var d Ding
	vd := reflect.ValueOf(&d).Elem()
	fmt.Println(vd)
	object_field := vd.FieldByName("Object")

	ptr := reflect.New(object_field.Type())
	p2 := ptr.Elem()
	ptr.Elem().Set(reflect.New(p2.Type().Elem()))

	setString(ptr.Elem().Interface(), "Lasse")
	fmt.Println(to)
}

func setString(v interface{}, value string) {
	val := reflect.ValueOf(v).Elem()
	val.FieldByName("Name").SetString(value)
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
