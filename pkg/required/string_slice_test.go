package required

import (
	"encoding/json"
	"testing"
)

type StringSliceChecker struct {
	Names StringSlice `json:"names"`
}

func TestNewStringSlice(t *testing.T) {
	v := NewStringSlice([]string{"one", "two", "three"})
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	var b StringSlice
	if err := Unmarshal(data, &b); err != nil {
		t.Fatal(err)
	}
	if len(b.value) != len(v.value) {
		t.Fatalf("%v != %v", len(b.value), len(v.value))
	}
}

func TestStringSliceValidation(t *testing.T) {
	tt := []struct {
		name   string
		json   string
		err    error
		assert func(v interface{}) bool
	}{
		{"valid strincg", `{"names":["Lasse", "Basse"]}`, nil, func(p interface{}) bool { return p.(StringSliceChecker).Names.Value()[0] == "Lasse" }},
		{"empty string", `{"names": []}`, ErrEmptyStringSlice, skipAssert},
		{"nil string", `{}`, ErrEmptyStringSlice, skipAssert},
	}

	for _, tf := range tt {
		t.Run(tf.name, func(t *testing.T) {
			jsonb := []byte(tf.json)
			var checker StringSliceChecker
			err := Unmarshal(jsonb, &checker)
			assertError(t, err, tf.err)

			if !tf.assert(checker) {
				t.Fatalf("Assertion Failed: %+v", checker)
			}
		})
	}
}
