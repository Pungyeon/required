package json

import (
	"encoding/json"
	"testing"
)

type TheThing struct{
	Ding
}

type Object struct {
	Name string `json:"name"`
}

type Ding struct {
	Ding int64 `json:"ding"`
	Dong string `json:"dong"`
	Float float64 `json:"float"`
	Object Object `json:"object"`
	Array []int `json:"array"`
	StringSlice []string `json:"string_slice"`
}
var sample = []byte(`{
		"ding": 1,
		"dong": "hello",
		"object": {
			"name": "lasse"
		},
		"float": 3.2,
		"array": [1, 2, 3],
		"string_slice": ["1", "2", "3"]
	}`)

func TestUnmarshal(t *testing.T) {
	var ding Ding
	if err := Unmarshal(sample, &ding); err != nil {
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

func BenchmarkStdUnmarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var ding Ding
		if err := json.Unmarshal(sample, &ding); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPkgUnmarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var ding Ding
		if err := Unmarshal(sample, &ding); err != nil {
			b.Fatal(err)
		}
	}
}
