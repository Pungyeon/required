package required

import (
	"encoding/json"
	"testing"
)

type Embedded struct {
	Bool         `json:"bool"`
	String       `json:"string"`
	Int          `json:"int"`
	Float32      `json:"float32"`
	Float64      `json:"float64"`
	IntSlice     `json:"int_slice"`
	ByteSlice    `json:"byte_slice"`
	BoolSlice    `json:"bool_slice"`
	Float32Slice `json:"float32_slice"`
	Float64Slice `json:"float64_slice"`
	StringSlice  `json:"string_slice"`
}

func TestEmbedded(t *testing.T) {
	jsonBytes := []byte(`{
		"bool": true,
		"string": "hello",
		"int": 64,
		"float32": 32.2,
		"float64": 64.4,
		"int_slice": [1],
		"byte_slice": [1],
		"bool_slice": [true],
		"float32_slice": [32.2],
		"float64_slice": [64.4],
		"string_slice": ["hello"]
	}`)

	var embed Embedded
	if err := json.Unmarshal(jsonBytes, &embed); err != nil {
		t.Fatal(err)
	}
	assertEmbedded(t, embed.Bool.Value() == true, "Bool")
	assertEmbedded(t, embed.String.Value() == "hello", "String")
	assertEmbedded(t, embed.Int.Value() == 64, "Int")
	assertEmbedded(t, embed.Float32.Value() == 32.2, "Float32")
	assertEmbedded(t, embed.Float64.Value() == 64.4, "Float64")
	assertEmbedded(t, embed.IntSlice.Value()[0] == 1, "IntSlice")
	assertEmbedded(t, embed.ByteSlice.Value()[0] == 1, "ByteSlice")
	assertEmbedded(t, embed.BoolSlice.Value()[0] == true, "BoolSlice")
	assertEmbedded(t, embed.Float32Slice.Value()[0] == 32.2, "Float32Slice")
	assertEmbedded(t, embed.Float64Slice.Value()[0] == 64.4, "Float64Slice")
	assertEmbedded(t, embed.StringSlice.Value()[0] == "hello", "StringSlice")
}

func assertEmbedded(t *testing.T, assertion bool, typeStr string) {
	t.Helper()
	if !assertion {
		t.Fatalf("Embedded %s parsing error", typeStr)
	}
}
