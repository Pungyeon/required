package unsafe

import (
	"fmt"
	"testing"
)

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

type Person struct {
	Name string `json:"name"`
	Age int `json:"age"`
	Twitter string `json:"twitter"`
}

func (p *Person) Hi() {
	fmt.Println("hello there")
}

