package json

import "testing"

type Ding struct {
	Dong string `json:"dong"`
}

func TestUnmarshal(t *testing.T) {
	json := []byte(`{
		"dong": "hello"
	}`)

	var ding Ding
	if err := Unmarshal(json, &ding); err != nil {
		t.Fatal(err)
	}

	if ding.Dong != "hello" {
		t.Fatalf("mismatch: (%s) != (%s)", ding.Dong, "hello")
	}
}