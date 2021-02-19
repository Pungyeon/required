package unsafe

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

type Human struct {
	Name string `json:"name"`
	Age int `json:"age"`
	Data uint8
}

func TestParseObject(t *testing.T) {
	var human Human
	if err := Unmarshal([]byte(`{"name": "lasse", "age": 30, "data": 9}`), &human); err != nil {
		t.Fatal(err)
	}
	if human.Name != "lasse" {
		t.Fatal("wrong name", human)
	}
	fmt.Println(human)
}

func TestMe(t *testing.T) {
	person := &Person{Name: "Lasse", Age: 23, Twitter: "ifndef_lmj"}
	fmt.Println(person)
	v := ValueOf(person)
	//v = Method(v, 0)
	Name(v, 0)
	fmt.Println(person)
}

func TestTags(t *testing.T) {
	person := Person{Name: "Lasse", Age: 23, Twitter: "ifndef_lmj"}
	v := ValueOf(person)
	tags, err := GetTags(v)
	if err != nil {
		t.Fatal(err)
	}
	if tags["name"].FieldName != "name" {
		t.Fatal(tags)
	}
}

func TestValueConversion(t *testing.T) {
	rval := reflect.ValueOf(&Person{
		Name:    "Lasse",
		Age:     30,
		Twitter: "ifndef_lmj",
	})

	uval := *(*Value)(unsafe.Pointer(&rval))
	elem := (*ptrType)(unsafe.Pointer(uval.typ)).elem
	fmt.Println(uval.Type(), uval.kind())
	fmt.Println(elem.Type(), elem.Kind())

}

type Person struct {
	Name string `json:"name"`
	Age int `json:"age"`
	Twitter string `json:"twitter"`
	Address Address
}

type Address struct {
	Street string
	Number int
}

func (p *Person) Hi() {
	fmt.Println("hello there")
}

