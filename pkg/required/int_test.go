package required

import (
	"encoding/json"
	"testing"
)

type Customer struct {
	ID Int `json:"age"`
}

func TestNewInt(t *testing.T) {
	v := NewInt(12)
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	var b Int
	if err := Unmarshal(data, &b); err != nil {
		t.Fatal(err)
	}
	if b.Value() != v.Value() {
		t.Fatalf("%v != %v", b.Value(), v.Value())
	}
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
		{"empty int", `{"name":""}`, ErrEmptyInt, skipAssert},
		{"nil int", `{}`, ErrEmptyInt, skipAssert},
	}

	for _, tf := range tt {
		t.Run(tf.name, func(t *testing.T) {
			jsonb := []byte(tf.json)
			var customer Customer
			err := Unmarshal(jsonb, &customer)
			assertError(t, err, tf.err)


			if !tf.assert(customer) {
				t.Fatalf("Assertion Failed: %+v", customer)
			}
		})
	}
}
