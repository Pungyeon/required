package required

import (
	"testing"
)

type Float32SliceChecker struct {
	Floats Float32Slice `json:"data"`
}

func TestFloat32SliceValidation(t *testing.T) {
	tt := []struct {
		name   string
		json   string
		err    error
		assert func(v interface{}) bool
	}{
		{"valid strincg", `{"data":[233,2,3,125]}`, nil, func(p interface{}) bool { return p.(Float32SliceChecker).Floats.Value()[0] == 233 }},
		{"empty string", `{"data": []}`, ErrEmptyFloatSlice, skipAssert},
		{"nil string", `{}`, ErrEmptyFloatSlice, skipAssert},
	}

	for _, tf := range tt {
		t.Run(tf.name, func(t *testing.T) {
			jsonb := []byte(tf.json)
			var checker Float32SliceChecker
			err := Unmarshal(jsonb, &checker)
			assertError(t, err, tf.err)


			if !tf.assert(checker) {
				t.Fatalf("Assertion Failed: %+v", checker)
			}
		})
	}
}

type Float64SliceChecker struct {
	Floats Float64Slice `json:"data"`
}

func TestFloat64SliceValidation(t *testing.T) {
	tt := []struct {
		name   string
		json   string
		err    error
		assert func(v interface{}) bool
	}{
		{"valid strincg", `{"data":[233,2,3,125]}`, nil, func(p interface{}) bool { return p.(Float64SliceChecker).Floats.Value()[0] == 233 }},
		{"empty string", `{"data": []}`, ErrEmptyFloatSlice, skipAssert},
		{"nil string", `{}`, ErrEmptyFloatSlice, skipAssert},
	}

	for _, tf := range tt {
		t.Run(tf.name, func(t *testing.T) {
			jsonb := []byte(tf.json)
			var checker Float64SliceChecker
			err := Unmarshal(jsonb, &checker)
			assertError(t, err, tf.err)


			if !tf.assert(checker) {
				t.Fatalf("Assertion Failed: %+v", checker)
			}
		})
	}
}
