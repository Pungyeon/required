package json

import (
	"encoding/json"
	"testing"
)

type MarshalObj struct {
	Name    string
	Integer int
	Float   float64
	Bool    bool
	Array   []int
}

func TestMarshalSupport(t *testing.T) {
	data, err := marshal("Name")
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `"Name"` {
		t.Fatal(string(data))
	}

	data, err = marshal(MarshalObj{
		Name:    "Lasse",
		Integer: 1,
		Bool:    true,
		Array:   []int{1, 2, 3},
		Float:   3.2,
	})
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"name":"Lasse","integer":1,"float":3.2,"bool":true,"array":[1,2,3]}` {
		t.Fatal(string(data))
	}
}

var objSample = MarshalObj{
	Name: "Lasse",
}

func BenchmarkMarshalStd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data, err := json.Marshal(objSample)
		if err != nil {
			b.Fatal(err)
		}
		_ = data
	}
}

func BenchmarkMarshalPkg(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data, err := marshal(objSample)
		if err != nil {
			b.Fatal(err)
		}
		_ = data
	}
}
