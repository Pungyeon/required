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
		{"empty bool", `{"name":"dingeling"}`, ErrEmpty, skipAssert},
		{"nil bool", `{}`, ErrEmpty, skipAssert},
	}

	for _, tf := range tt {
		t.Run(tf.name, func(t *testing.T) {
			jsonb := []byte(tf.json)
			var reqBool RequiredBool
			if err := Unmarshal(jsonb, &reqBool); err != tf.err {
				t.Fatal(err)
			}

			if !tf.assert(reqBool) {
				t.Fatalf("Assertion Failed: %+v", reqBool)
			}
		})
	}
}
