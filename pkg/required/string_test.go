package required

import (
	"encoding/json"
	"testing"
)

type Person struct {
	Name String `json:"name"`
	Age  int64  `json:"age"`
}

func TestNewString(t *testing.T) {
	v := NewString("string")
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	var b String
	if err := Unmarshal(data, &b); err != nil {
		t.Fatal(err)
	}
	if b.Value() != v.Value() {
		t.Fatalf("%v != %v", b.Value(), v.Value())
	}
}

func TestStringValidation(t *testing.T) {
	tt := []struct {
		name   string
		json   string
		err    error
		assert func(v interface{}) bool
	}{
		{"valid string", `{"name":"Lasse"}`, nil, func(p interface{}) bool { return p.(Person).Name.Value() == "Lasse" }},
		{"empty string", `{"name":""}`, ErrEmptyString, skipAssert},
		{"nil string", `{}`, ErrEmptyString, skipAssert},
	}

	for _, tf := range tt {
		t.Run(tf.name, func(t *testing.T) {
			jsonb := []byte(tf.json)
			var person Person
			err := Unmarshal(jsonb, &person)
			assertError(t, err, tf.err)


			if !tf.assert(person) {
				t.Fatalf("Assertion Failed: %+v", person)
			}
		})
	}
}
