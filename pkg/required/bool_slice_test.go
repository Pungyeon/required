package required

import (
	"encoding/json"
	"testing"
)

type BoolSliceChecker struct {
	Bools BoolSlice `json:"data"`
}

func TestNewBoolSlice(t *testing.T) {
	v := NewBoolSlice([]bool{true, false, true})
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	var b BoolSlice
	if err := Unmarshal(data, &b); err != nil {
		t.Fatal(err)
	}
	if len(b.value) != 3 {
		t.Fatalf("%v != %v", len(b.value), 3)
	}
}

func TestBoolSliceValidation(t *testing.T) {
	tt := []struct {
		name   string
		json   string
		err    error
		assert func(v interface{}) bool
	}{
		{"valid bool slice", `{"data":[true, false, true]}`, nil, func(p interface{}) bool { return p.(BoolSliceChecker).Bools.Value()[0] == true }},
		{"empty bool slice", `{"data": []}`, ErrEmptyBoolSlice, skipAssert},
		{"nil bool slice", `{}`, ErrEmptyBoolSlice, skipAssert},
	}

	for _, tf := range tt {
		t.Run(tf.name, func(t *testing.T) {
			jsonb := []byte(tf.json)
			var checker BoolSliceChecker
			err := Unmarshal(jsonb, &checker)
			assertError(t, err, tf.err)

			if !tf.assert(checker) {
				t.Fatalf("Assertion Failed: %+v", checker)
			}
		})
	}
}
