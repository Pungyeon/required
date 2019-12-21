package required

import (
	"testing"
)

type BoolSliceChecker struct {
	Bools BoolSlice `json:"data"`
}

func TestBoolSliceValidation(t *testing.T) {
	tt := []struct {
		name   string
		json   string
		err    error
		assert func(v interface{}) bool
	}{
		{"valid strincg", `{"data":[true, false, true]}`, nil, func(p interface{}) bool { return p.(BoolSliceChecker).Bools.Value()[0] == true }},
		{"empty string", `{"data": []}`, ErrEmpty, skipAssert},
		{"nil string", `{}`, ErrEmpty, skipAssert},
	}

	for _, tf := range tt {
		t.Run(tf.name, func(t *testing.T) {
			jsonb := []byte(tf.json)
			var checker BoolSliceChecker
			if err := Unmarshal(jsonb, &checker); err != tf.err {
				t.Fatal(err)
			}

			if !tf.assert(checker) {
				t.Fatalf("Assertion Failed: %+v", checker)
			}
		})
	}
}
