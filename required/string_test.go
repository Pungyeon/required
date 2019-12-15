package required

import (
	"testing"
)

type Person struct {
	Name String `json:"name"`
	Age  int64  `json:"age"`
}

func skipAssert(p Person) bool {
	return true
}

func TestStringValidation(t *testing.T) {
	tt := []struct {
		name   string
		json   string
		err    error
		assert func(Person) bool
	}{
		{"valid strincg", `{"name":"Lasse"}`, nil, func(p Person) bool { return p.Name.Value() == "Lasse" }},
		{"empty string", `{"name":""}`, ErrStringEmpty, skipAssert},
		{"nil string", `{}`, ErrStringEmpty, skipAssert},
	}

	for _, tf := range tt {
		t.Run(tf.name, func(t *testing.T) {
			jsonb := []byte(tf.json)
			var person Person
			if err := Unmarshal(jsonb, &person); err != tf.err {
				t.Fatal(err)
			}

			if !tf.assert(person) {
				t.Fatalf("Assertion Failed: %+v", person)
			}
		})
	}
}
