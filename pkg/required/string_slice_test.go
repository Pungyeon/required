package required

import (
	"testing"
)

type StringSliceChecker struct {
	Names StringSlice `json:"names"`
}

func TestStringSliceValidation(t *testing.T) {
	tt := []struct {
		name   string
		json   string
		err    error
		assert func(v interface{}) bool
	}{
		{"valid strincg", `{"names":["Lasse", "Basse"]}`, nil, func(p interface{}) bool { return p.(StringSliceChecker).Names.Value()[0] == "Lasse" }},
		{"empty string", `{"names": []}`, ErrEmpty, skipAssert},
		{"nil string", `{}`, ErrEmpty, skipAssert},
	}

	for _, tf := range tt {
		t.Run(tf.name, func(t *testing.T) {
			jsonb := []byte(tf.json)
			var checker StringSliceChecker
			if err := Unmarshal(jsonb, &checker); err != tf.err {
				t.Fatal(err)
			}

			if !tf.assert(checker) {
				t.Fatalf("Assertion Failed: %+v", checker)
			}
		})
	}
}
