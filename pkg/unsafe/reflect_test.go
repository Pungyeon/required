package unsafe

import (
	"fmt"
	"github.com/Pungyeon/required/pkg/json"
	"reflect"
	"testing"
	"unsafe"
)

type Human struct {
	Name string `json:"name"`
	Age int `json:"age"`
	Data uint8
	Address Address
	Array []int
	Alive bool `json:"alive"`
}

var example = []byte(`{
	"name": "lasse",
	"age": 30,
	"data": 9,
	"address": {
		"street": "privet drive"
	},
	"array": [1, 2, 3, 4, 5, 6],
	"alive": true
}`)

func BenchmarkParseObjectUnsafe(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var human Human
		if err := Unmarshal(example, &human); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseObject(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var human Human
		if err := json.Unmarshal(example, &human); err != nil {
			b.Fatal(err)
		}
	}
}

func TestParseMap(t *testing.T) {
}

func TestParseArray(t *testing.T) {
	arr := []int{
		1, 2, 3, 4, 5,
	}
	val := ValueOf(&arr)
	fmt.Println(val.typ.Kind())

	var tt *sliceType
	if val.typ.Kind() == reflect.Ptr {
		ptr := (*ptrType)(unsafe.Pointer(val.typ))
		tt = (*sliceType)(unsafe.Pointer(ptr.elem))
	} else {
		tt = (*sliceType)(unsafe.Pointer(val.typ))
	}

	s := (*Slice)(val.ptr)
	fmt.Println(s.Len)

	for i := 0; i < s.Len; i++ {
		v := add(s.Data, uintptr(i)*tt.elem.size)
		fmt.Println(*(*int)(v))
	}

	*(*[]*rtype)(val.ptr) = make([]*rtype, 2)

	v := getSliceIndex(s, tt, 1)
	*(*int)(v) = 2

	v = getSliceIndex(s, tt, 0)
	*(*int)(v) = 1
	//reflect.MakeSlice(reflect.Int, 1).Index()
	//*(*[]int)(val.ptr) = []int{1, 2, 3}
	fmt.Println(arr)
}


func TestParsePrimitive(t *testing.T) {
	var i int
	if err := Unmarshal([]byte(`32`), &i); err != nil {
		t.Fatal(err)
	}
	if i != 32 {
		t.Fatal(i)
	}

	var s string
	if err := Unmarshal([]byte(`"dingeling"`), &s); err != nil {
		t.Fatal(err)
	}
	if s != "dingeling" {
		t.Fatal(s)
	}

	var float float32
	if err := Unmarshal([]byte(`3.2`), &float); err != nil {
		t.Fatal(err)
	}
	if float != 3.2 {
		t.Fatal(float)
	}

	var b bool
	if err := Unmarshal([]byte(`true`), &b); err != nil {
		t.Fatal(err)
	}
	if !b {
		t.Fatal(b)
	}
}

func TestParseObject(t *testing.T) {
	var human Human
	if err := Unmarshal(example, &human); err != nil {
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

