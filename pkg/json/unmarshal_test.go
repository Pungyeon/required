package json

import "testing"

type Object struct {
	Name string `json:"name"`
}

type Ding struct {
	Ding int64 `json:"ding"`
	Dong string `json:"dong"`
	Float float64 `json:"float"`
	Object Object `json:"object"`
}

func TestUnmarshal(t *testing.T) {
	json := []byte(`{
		"ding": 1,
		"dong": "hello",
		"float": 3.2,
		"object": {
			"name": "lasse"
		}
	}`)

	var ding Ding
	if err := Unmarshal(json, &ding); err != nil {
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
}