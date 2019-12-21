package required

import (
	"testing"
)

type Customer struct {
	ID   Int `json:"age"`
	Bool `json:"name"`
}

func skipAssert(v interface{}) bool {
	return true
}

func TestIntValidation(t *testing.T) {
	tt := []struct {
		name   string
		json   string
		err    error
		assert func(v interface{}) bool
	}{
		{"valid int", `{"age": 29}`, nil, func(c interface{}) bool { return c.(Customer).ID.Value() == 29 }},
		{"empty int", `{"name":""}`, ErrEmpty, skipAssert},
		{"nil int", `{}`, ErrEmpty, skipAssert},
	}

	for _, tf := range tt {
		t.Run(tf.name, func(t *testing.T) {
			jsonb := []byte(tf.json)
			var customer Customer
			if err := Unmarshal(jsonb, &customer); err != tf.err {
				t.Fatal(err)
			}

			if !tf.assert(customer) {
				t.Fatalf("Assertion Failed: %+v", customer)
			}
		})
	}
}
