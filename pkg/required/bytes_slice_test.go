package required

import (
	"testing"
)

type ByteSliceChecker struct {
	Bytes ByteSlice `json:"data"`
}

func TestByteSliceValidation(t *testing.T) {
	tt := []struct {
		name   string
		json   string
		err    error
		assert func(v interface{}) bool
	}{
		{"valid strincg", `{"data":[233,2,3,125]}`, nil, func(p interface{}) bool { return p.(ByteSliceChecker).Bytes.Value()[0] == 233 }},
		{"empty string", `{"data": []}`, ErrEmptyByteSlice, skipAssert},
		{"nil string", `{}`, ErrEmptyByteSlice, skipAssert},
	}

	for _, tf := range tt {
		t.Run(tf.name, func(t *testing.T) {
			jsonb := []byte(tf.json)
			var checker ByteSliceChecker
			err := Unmarshal(jsonb, &checker)
			assertError(t, err, tf.err)


			if !tf.assert(checker) {
				t.Fatalf("Assertion Failed: %+v", checker)
			}
		})
	}
}
