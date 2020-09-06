package required

import (
	"encoding/json"
	"testing"
)


type FloatChecker struct {
	Thing Float `json:"thing"`
}

func TestNewFloat(t *testing.T) {
	v := NewFloat(32.2)
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	var b Float
	if err := Unmarshal(data, &b); err != nil {
		t.Fatal(err)
	}
	if b.Value() != 32.2 {
		t.Fatalf("%v != %v", b.Value(), 32.2)
	}
}


func TestFloatValidation(t *testing.T) {
	tt := []struct {
		name   string
		json   string
		err    error
		assert func(v interface{}) bool
	}{
		{"valid float", `{"thing": 52.3}`, nil, func(c interface{}) bool { return c.(FloatChecker).Thing.Value() == 52.3 }},
		{"valid int", `{"thing": 52}`, nil, func(c interface{}) bool { return c.(FloatChecker).Thing.Value() == 52 }},
		{"empty", `{"name":""}`, ErrEmptyFloat, skipAssert},
		{"nil", `{}`, ErrEmptyFloat, skipAssert},
	}

	for _, tf := range tt {
		t.Run(tf.name, func(t *testing.T) {
			jsonb := []byte(tf.json)
			var checker FloatChecker
			err := Unmarshal(jsonb, &checker)
			assertError(t, err, tf.err)


			if !tf.assert(checker) {
				t.Fatalf("Assertion Failed: %+v", checker)
			}
		})
	}
}
