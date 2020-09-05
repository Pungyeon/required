package required

import (
	"testing"
)

type IntSliceChecker struct {
	Ints IntSlice `json:"data"`
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
