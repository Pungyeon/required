package json

import (
	"reflect"
	"testing"
)

type TagTest struct {
	Name     string
	DingDong string
}

func TestTagFormatting(t *testing.T) {
	v := reflect.ValueOf(TagTest{
		Name:     "lasse",
		DingDong: "ding_dong",
	})
	tags := getFieldTags(v)
	if _, ok := tags["name"]; !ok {
		t.Fatal(`could not find "name" tag`)
	}

	if _, ok := tags["ding_dong"]; !ok {
		t.Fatal(`could not find "ding_dong" tag`)
	}
}

func TestToSnakeCase(t *testing.T) {
	camel := "DingDong"
	if toSnakeCase(camel) != "ding_dong" {
		t.Fatal("oh dear", toSnakeCase(camel))
	}
}
