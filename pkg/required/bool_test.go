package required

import (
	"testing"
)

type RequiredBool struct {
	Active Bool `json:"active"`
	Name   string
}

func TestBoolValidation(t *testing.T) {
	tt := []struct {
		name   string
		json   string
		err    error
		assert func(v interface{}) bool
	}{
		{"valid bool", `{"active": true}`, nil, func(v interface{}) bool { return v.(RequiredBool).Active.Value() }},
		{"empty bool", `{"name":"dingeling"}`, ErrEmptyBool, skipAssert},
		{"nil bool", `{}`, ErrEmptyBool, skipAssert},
	}

	for _, tf := range tt {
		t.Run(tf.name, func(t *testing.T) {
			jsonb := []byte(tf.json)
			var reqBool RequiredBool
			err := Unmarshal(jsonb, &reqBool)
			assertError(t, err, tf.err)

			if !tf.assert(reqBool) {
				t.Fatalf("Assertion Failed: %+v", reqBool)
			}
		})
	}
}
