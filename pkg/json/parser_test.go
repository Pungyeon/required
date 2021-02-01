package json

import (
	"bytes"
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"testing"

	"github.com/Pungyeon/required/pkg/token"

	"github.com/Pungyeon/required/pkg/lexer"
	"github.com/Pungyeon/required/pkg/structtag"
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

func LexString(t *testing.T, input string) *lexer.Lexer {
	return lexer.NewLexer([]byte(input))
}

type PrivateFields struct {
	ding string
	dong TestObject
}

func TestPrivateFields(t *testing.T) {
	var pf PrivateFields
	if err := Parse(LexString(t, `{"ding": "dingeling", "dong": {"name": "lasse"}}`), &pf); err != nil {
		t.Fatal(err)
	}
}

func TestParserSimple(t *testing.T) {
	var obj TestObject
	if err := Parse(LexString(t, `{"name": "lasse"}`), &obj); err != nil {
		t.Fatal(err)
	}
	if obj.Name != "lasse" {
		t.Fatal("not lasse:", obj.Name)
	}
}

func TestParseSample(t *testing.T) {
	var ding Ding
	if err := Parse(LexString(t, sample), &ding); err != nil {
		t.Fatal(err)
	}
	if !(ding.Ding == 1 &&
		ding.Dong == "hello" &&
		ding.Boolean == true &&
		ding.Object.Name == "lasse" &&
		ding.Array[2] == 3 &&
		ding.StringSlice[2] == "3" &&
		ding.MultiDimension[1][2] == 6 &&
		ding.ObjectArray[1].Name == "basse" &&
		ding.MapObject["lumber"] == 13 &&
		ding.Float == 3.2) {
		t.Fatal(ding)
	}
}

func TestParsePrimitive(t *testing.T) {
	var v uint64
	if err := Parse(LexString(t, `1`), &v); err != nil {
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
	tokens := LexString(t, `{"array": [1, 2, 3, 4]}`)

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

	if obj.Array[3] != 4 {
		t.Fatal("expected 4:", obj.Array[3])
	}
}

func TestParseArray(t *testing.T) {
	tokens := LexString(t, "[1, 2, 3, 4]")
	//if tokens.Join(";") != "[;1;,;2;,;3;,;4;]" {
	//	t.Fatal("oh no", tokens.Join(";"))
	//}

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
	tokens := LexString(t, "[1.1, 2.2, 3.3, 4.4]")
	//if tokens.Join(";") != "[;1.1;,;2.2;,;3.3;,;4.4;]" {
	//	t.Fatal("oh no", tokens.Join(";"))
	//}

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
	tokens := LexString(t, `[
	[1, 2, 3],
	[4, 5, 6]
]`)
	//if tokens.Join(";") != "[;[;1;,;2;,;3;];,;[;4;,;5;,;6;];]" {
	//	t.Fatal("oh no", tokens.Join(";"))
	//}

	var obj [][]int
	if err := Parse(tokens, &obj); err != nil {
		t.Fatal(err)
	}

	if len(obj) != 2 {
		t.Fatal("length of object:", len(obj), obj)
	}

	if len(obj[0]) != 3 {
		t.Fatal("length of object inner:", len(obj[0]))
	}

	if obj[1][0] != 4 {
		t.Fatal("omg it's not 4:", obj[1][0])
	}
}

func TestParseMultiStringArray(t *testing.T) {
	tokens := LexString(t, `[
	["1", "2", "3"],
	["4", "5", "6"]
]`)
	//if tokens.Join(";") != "[;[;1;,;2;,;3;];,;[;4;,;5;,;6;];]" {
	//	t.Fatal("oh no", tokens.Join(";"))
	//}

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
	tokens := LexString(t, `[
	{
		"name": "lasse"
	},
	{
		"name": "basse"
	}
]`)
	//if tokens.Join(";") != "[;{;name;:;lasse;};,;{;name;:;basse;};]" {
	//	t.Fatal("oh no", tokens.Join(";"))
	//}

	var obj []TestObject
	if err := Parse(tokens, &obj); err != nil {
		t.Fatal(err)
	}

	if len(obj) != 2 {
		t.Fatal("length of object:", len(obj), obj)
	}

	if obj[1].Name != "basse" {
		t.Fatal("omg it's not basse:", obj[1].Name)
	}
}

func TestMapStringIntUnmarshal(t *testing.T) {
	tokens := LexString(t, `{"number": 1, "lumber": 13}`)
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
	tokens := LexString(t, `{"number": "1", "lumber": "13"}`)
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

func testParse(t *testing.T, lexer *lexer.Lexer, v interface{}) {
	if err := Parse(lexer, v); err != nil {
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
			testParse(t, LexString(t, `"lasse"`), &v)
			return v == "lasse"
		}},
		{name: "int64", check: func() bool {
			var v int64
			testParse(t, LexString(t, `234`), &v)
			return v == 234
		}},
		{name: "int32", check: func() bool {
			var v int32
			testParse(t, LexString(t, `234`), &v)
			return v == 234
		}},
		{name: "int", check: func() bool {
			var v int
			testParse(t, LexString(t, `234`), &v)
			return v == 234
		}},
		{name: "float64", check: func() bool {
			var v float64
			testParse(t, LexString(t, `42.2`), &v)
			return v == 42.2
		}},
		{name: "float32", check: func() bool {
			var v float32
			testParse(t, LexString(t, `42.2`), &v)
			return v == 42.2
		}},
		{name: "test_object", check: func() bool {
			var v TestObject
			testParse(t, LexString(t, `{"name": "lasse"}`), &v)
			return v.Name == "lasse"
		}},
		{name: "string_array", check: func() bool {
			var v []string
			testParse(t, LexString(t, `["name", "lasse"]`), &v)
			return v[1] == "lasse"
		}},
		{name: "interface_string", check: func() bool {
			var v interface{}
			testParse(t, LexString(t, `"lasse"`), &v)
			return v.(string) == "lasse"
		}},
		{name: "interface_object", check: func() bool {
			var v interface{}
			testParse(t, LexString(t, `{"name": "lasse"}`), &v)
			return v.(map[string]interface{})["name"] == "lasse"
		}},
		{name: "interface_array", check: func() bool {
			var v []interface{}
			testParse(t, LexString(t, `["name", "lasse"]`), &v)
			return v != nil &&
				v[0].(string) == "name"
		}},
		{name: "ding_object", check: func() bool {
			var ding Ding
			testParse(t, LexString(t, sample), &ding)
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
	tokens := LexString(t, `{
		"object": {
			"name": "lasse"
		}
	}`)

	var ding Ding
	if err := Parse(tokens, &ding); err != nil {
		t.Fatal(err)
	}

	if ding.Object.Name != "lasse" {
		t.Fatal("oh no")
	}
}

func TestParseInterfaceString(t *testing.T) {
	tokens := LexString(t, `"lasse"`)

	var ding interface{}
	if err := Parse(tokens, &ding); err != nil {
		t.Fatal(err)
	}

	if ding.(string) != "lasse" {
		t.Fatal("oh no")
	}
}

func TestMapFollowedBy(t *testing.T) {
	tokens := LexString(t, `{
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

type CustomRequiredEmail string

var errEmailRequired = errors.New("email field required")

func (email CustomRequiredEmail) IsValueValid() error {
	matched, err := regexp.MatchString(`.+@.+\..+`, string(email))
	if err != nil {
		return err
	}
	if !matched {
		return errEmailRequired
	}
	return nil
}

func TestRequiredFields(t *testing.T) {
	type RequiredBoi struct {
		Name string `json:"name, required"`
	}

	var r RequiredBoi
	if err := Parse(LexString(t, `{}`), &r); !structtag.IsRequiredErr(err) {
		t.Fatal("no required error, or unexpected error returned:", err)
	}

	if err := Parse(LexString(t, `{"name": "lasse"}`), &r); err != nil {
		t.Fatal(err)
	}
	type TestUser struct {
		Email CustomRequiredEmail
	}

	var invalidEmail TestUser
	if err := Parse(LexString(t, `{"email": "dingeling.dk"}`), &invalidEmail); err != errEmailRequired {
		t.Fatal("no required error, or unexpected error returned:", err)
	}

	var validEmail TestUser
	if err := Parse(LexString(t, `{"email": "lasse@jakobsen.dev"}`), &validEmail); err != nil {
		t.Fatal("no required error, or unexpected error returned:", err)
	}
}

func TestNullSupport(t *testing.T) {
	var d Ding
	if err := Parse(LexString(t, `{"object": null}`), &d); err != nil {
		t.Fatal(err)
	}
	if d.Object != nil {
		t.Fatal("object not nil")
	}

	var to TestObject
	if err := Parse(LexString(t, `{"name": null}`), &to); err != nil {
		t.Fatal(err)
	}
	if to.Name != "" {
		t.Fatal("name not nothing:", to.Name)
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

type IntString struct {
	Value int
}

func (is *IntString) UnmarshalJSON(data []byte) error {
	type Anon struct {
		Value string
	}
	var v Anon
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	val, err := strconv.ParseInt(v.Value, 10, 64)
	if err != nil {
		return err
	}
	is.Value = int(val)
	return nil
}

var dateSample = []byte(`{
	"value": "123"
}`)

func TestCustomMarshaler(t *testing.T) {
	var d IntString
	if err := Parse(lexer.NewLexer(dateSample), &d); err != nil {
		t.Fatal(err)
	}

	if d.Value != 123 {
		t.Fatal(d.Value)
	}
}

func TestDecoder(t *testing.T) {
	r := bytes.NewReader([]byte(sample))
	var ding Ding
	if err := NewDecoder(r).Decode(&ding); err != nil {
		t.Fatal(err)
	}
	defer recovery(t, ding)
	if !(ding.Ding == 1 &&
		ding.Dong == "hello" &&
		ding.Boolean == true &&
		ding.Object.Name == "lasse" &&
		len(ding.Array) == 3) {
		//ding.StringSlice[2] == "3" &&
		//ding.MultiDimension[1][2] == 6 &&
		//ding.ObjectArray[1].Name == "basse" &&
		//ding.MapObject["lumber"] == 13 &&
		//ding.Float == 3.2) {
		t.Fatal(ding)
	}
}

func TestInvalidJSON(t *testing.T) {
	tt := []struct {
		name string
		json string
		err  error
	}{
		{"early comma", `{"name",}`, token.ErrInvalidJSON},
		{"wrong value", `{"ding": "this should be an int64"`, token.ErrInvalidValue},
		{"unfinished inner array", `{"multi_dimension": [[]`, token.ErrUnmatchedBrace},
		{"unfinished object", `{"name"}`, token.ErrInvalidJSON},
		{"not a bool", `{"boolean": "yes"}`, token.ErrInvalidValue},
	}

	for _, tf := range tt {
		t.Run(tf.name, func(t *testing.T) {
			var ding Ding
			if err := Parse(LexString(t, tf.json), &ding); !errors.Is(err, tf.err) {
				t.Fatal(err)
			}
		})
	}
}

func TestMapArray(t *testing.T) {
	type MA struct {
		Array []map[string]int
	}

	var ma MA
	if err := Parse(LexString(t, `{ "array": [{ "name": 3 }] }`), &ma); err != nil {
		t.Fatal(err)
	}
	if ma.Array[0]["name"] != 3 {
		t.Fatal(ma)
	}
}

func TestMapArrayInterface(t *testing.T) {
	type MA struct {
		Array []map[string]int
	}

	var ma interface{}
	if err := Parse(LexString(t, `{ "array": [{ "name": 3 }] }`), &ma); err != nil {
		t.Fatal(err)
	}
	defer func() {
		r := recover()
		if r != nil {
			t.Fatal(r, ma)
		}
	}()

	v := ma.(map[string]interface{})["array"].([]interface{})[0].(map[string]interface{})["name"].(int)
	if v != 3 {
		t.Fatal(ma)
	}
}

func recovery(t *testing.T, v interface{}) {
	t.Helper()
	if r := recover(); r != nil {
		t.Error("recovered from panic:")
		t.Error(v)
		t.Fatal(r)
	}
}
