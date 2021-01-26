package json

import (
	"encoding/json"
	"errors"
	"testing"
)

type MarshalObj struct {
	Name      string
	Integer   int
	Float     float64
	Bool      bool
	Array     []int `json:"slice,required"`
	Map       map[int]string
	Struct    SmallObj
	Pointer   *SmallObj
	Interface interface{}

	// unsupported types
	channel  chan string
	function func()
	complex  complex128
	uptr     uintptr
}

type SmallObj struct {
	Name string
}

func TestMarshalUnsupportedType(t *testing.T) {
	var cmplx complex128
	var uptr uintptr
	tt := []struct {
		Name  string
		Value interface{}
		Error error
	}{
		{"channel", make(chan bool), ErrUnsupportedType},
		{"function", func() error { return nil }, ErrUnsupportedType},
		{"complex", cmplx, ErrUnsupportedType},
		{"uintptr", uptr, ErrUnsupportedType},
	}

	for _, tf := range tt {
		if _, err := Marshal(tf.Value); !errors.Is(err, ErrUnsupportedType) {
			t.Fatal(err)
		}
	}
}

func TestMarshalNullSupport(t *testing.T) {
	data, err := marshal(MarshalObj{})
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != `{"name":"","integer":0,"float":0,"bool":false,"slice":null,"map":null,"struct":{"name":""},"pointer":null,"interface":null}` {
		t.Fatal(string(data))
	}
}

func TestMarshalSupport(t *testing.T) {
	data, err := marshal("Name")
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `"Name"` {
		t.Fatal(string(data))
	}

	data, err = marshal(obj)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != `{"name":"Lasse","integer":1,"float":3.2,"bool":true,"slice":[1,2,3],"map":{"1":"hello","2":"goodbye"},"struct":{"name":"lasse"},"pointer":{"name":"pointer"},"interface":{"name":"interface"}}` {
		t.Fatal(string(data))
	}
}

var obj = MarshalObj{
	Name:    "Lasse",
	Integer: 1,
	Bool:    true,
	Array:   []int{1, 2, 3},
	Float:   3.2,
	Map: map[int]string{
		1: "hello",
		2: "goodbye",
	},
	Struct: SmallObj{
		Name: "lasse",
	},
	Pointer: &SmallObj{
		Name: "pointer",
	},
	Interface: &SmallObj{
		Name: "interface",
	},
}

func BenchmarkMarshalStd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data, err := json.Marshal(obj)
		if err != nil {
			b.Fatal(err)
		}
		_ = data
	}
}

func BenchmarkMarshalPkg(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data, err := marshal(obj)
		if err != nil {
			b.Fatal(err)
		}
		_ = data
	}
}
