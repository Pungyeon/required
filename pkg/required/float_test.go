package required

import (
	"testing"
)

type Float32Checker struct {
	Thing Float32 `json:"thing"`
}

func TestFloat32Validation(t *testing.T) {
	tt := []struct {
		name   string
		json   string
		err    error
		assert func(v interface{}) bool
	}{
		{"valid float32", `{"thing": 52.3}`, nil, func(c interface{}) bool { return c.(Float32Checker).Thing.Value() == 52.3 }},
		{"valid int", `{"thing": 52}`, nil, func(c interface{}) bool { return c.(Float32Checker).Thing.Value() == 52 }},
		{"empty", `{"name":""}`, ErrEmptyFloat, skipAssert},
		{"nil", `{}`, ErrEmptyFloat, skipAssert},
	}

	for _, tf := range tt {
		t.Run(tf.name, func(t *testing.T) {
			jsonb := []byte(tf.json)
			var checker Float32Checker
			err := Unmarshal(jsonb, &checker)
			assertError(t, err, tf.err)


			if !tf.assert(checker) {
				t.Fatalf("Assertion Failed: %+v", checker)
			}
		})
	}
}

type Float64Checker struct {
	Thing Float64 `json:"thing"`
}

func TestFloat64Validation(t *testing.T) {
	tt := []struct {
		name   string
		json   string
		err    error
		assert func(v interface{}) bool
	}{
		{"valid float", `{"thing": 52.3}`, nil, func(c interface{}) bool { return c.(Float64Checker).Thing.Value() == 52.3 }},
		{"valid int", `{"thing": 52}`, nil, func(c interface{}) bool { return c.(Float64Checker).Thing.Value() == 52 }},
		{"empty", `{"name":""}`, ErrEmptyFloat, skipAssert},
		{"nil", `{}`, ErrEmptyFloat, skipAssert},
	}

	for _, tf := range tt {
		t.Run(tf.name, func(t *testing.T) {
			jsonb := []byte(tf.json)
			var checker Float64Checker
			err := Unmarshal(jsonb, &checker)
			assertError(t, err, tf.err)


			if !tf.assert(checker) {
				t.Fatalf("Assertion Failed: %+v", checker)
			}
		})
	}
}
