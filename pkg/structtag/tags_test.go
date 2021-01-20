package structtag

import (
	"reflect"
	"testing"
)

type TagTest struct {
	Name     string
	DingDong string
	Age      int64 `json:"age,RequiredInterface"`
}

func TestTagFormatting(t *testing.T) {
	v := reflect.ValueOf(TagTest{
		Name:     "lasse",
		DingDong: "ding_dong",
		Age:      30,
	})
	tags, err := FromValue(v)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := tags["name"]; !ok {
		t.Fatal(`could not find "name" tag`)
	}

	dingDong, ok := tags["ding_dong"]
	if !ok {
		t.Fatal(`could not find "ding_dong" tag`)
	}
	if dingDong.FieldIndex != 1 {
		t.Fatal("unexpected field_index for ding_dong:", dingDong.FieldIndex)
	}

	age, ok := tags["age"]
	if !ok {
		t.Fatal(`could not find "age" tag`)
	}
	if age.FieldIndex != 2 {
		t.Fatal("unexpected field_index for age:", age.FieldIndex)
	}
}

func TestToSnakeCase(t *testing.T) {
	camel := "DingDong"
	if toSnakeCase(camel) != "ding_dong" {
		t.Fatal("oh dear", toSnakeCase(camel))
	}
}
