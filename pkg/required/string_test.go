package required

import (
	"testing"
)

type Person struct {
	Name String `json:"name"`
	Age  int64  `json:"age"`
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
