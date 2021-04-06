package required

import (
	"encoding/json"
	"testing"
)

type IntSliceChecker struct {
	Ints IntSlice `json:"data"`
}

func TestNewIntSlice(t *testing.T) {
	v := NewIntSlice([]int{1, 2, 3})
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	var b IntSlice
	if err := Unmarshal(data, &b); err != nil {
		t.Fatal(err)
	}
	if len(b.value) != len(v.value) {
		t.Fatalf("%v != %v", len(b.value), len(v.value))
	}
}

func TestIntSliceValidation(t *testing.T) {
	tt := []struct {
		name   string
		json   string
		err    error
		assert func(v interface{}) bool
	}{
		{"valid strincg", `{"data":[233,2,3,125]}`, nil, func(p interface{}) bool { return p.(IntSliceChecker).Ints.Value()[0] == 233 }},
		{"empty string", `{"data": []}`, ErrEmptyIntSlice, skipAssert},
		{"nil string", `{}`, ErrEmptyIntSlice, skipAssert},
	}

	for _, tf := range tt {
		t.Run(tf.name, func(t *testing.T) {
			jsonb := []byte(tf.json)
			var checker IntSliceChecker
			err := Unmarshal(jsonb, &checker)
			assertError(t, err, tf.err)

			if !tf.assert(checker) {
				t.Fatalf("Assertion Failed: %+v", checker)
			}
		})
	}
}
