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
		{"valid int", `{"thing": 52.3}`, nil, func(c interface{}) bool { return c.(Float32Checker).Thing.Value() == 52.3 }},
		{"empty int", `{"name":""}`, ErrEmpty, skipAssert},
		{"nil int", `{}`, ErrEmpty, skipAssert},
	}

	for _, tf := range tt {
		t.Run(tf.name, func(t *testing.T) {
			jsonb := []byte(tf.json)
			var checker Float32Checker
			if err := Unmarshal(jsonb, &checker); err != tf.err {
				t.Fatal(err)
			}

			if !tf.assert(checker) {
				t.Fatalf("Assertion Failed: %+v", checker)
			}
		})
	}
}
