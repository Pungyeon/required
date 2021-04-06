package required

import (
	"encoding/json"
	"testing"
)

type FloatSliceChecker struct {
	Floats FloatSlice `json:"data"`
}

func TestNewFloatSlice(t *testing.T) {
	v := NewFloatSlice([]float64{32.2, 64.4, 128.8})
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	var b FloatSlice
	if err := Unmarshal(data, &b); err != nil {
		t.Fatal(err)
	}
	if len(b.value) != len(v.value) {
		t.Fatalf("%v != %v", len(b.value), len(v.value))
	}
}

func TestFloatSliceValidation(t *testing.T) {
	tt := []struct {
		name   string
		json   string
		err    error
		assert func(v interface{}) bool
	}{
		{"valid strincg", `{"data":[233,2,3,125]}`, nil, func(p interface{}) bool { return p.(FloatSliceChecker).Floats.Value()[0] == 233 }},
		{"empty string", `{"data": []}`, ErrEmptyFloatSlice, skipAssert},
		{"nil string", `{}`, ErrEmptyFloatSlice, skipAssert},
	}

	for _, tf := range tt {
		t.Run(tf.name, func(t *testing.T) {
			jsonb := []byte(tf.json)
			var checker FloatSliceChecker
			err := Unmarshal(jsonb, &checker)
			assertError(t, err, tf.err)

			if !tf.assert(checker) {
				t.Fatalf("Assertion Failed: %+v", checker)
			}
		})
	}
}
