package required

import (
	"encoding/json"
	"testing"
)

type RequiredBool struct {
	Active Bool `json:"active"`
	Name   string
}

func TestNewBool(t *testing.T) {
	v := NewBool(true)
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	var b Bool
	if err := Unmarshal(data, &b); err != nil {
		t.Fatal(err)
	}
	if b.Value() != v.Value() {
		t.Fatalf("%v != %v", b.Value(), v.Value())
	}
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
