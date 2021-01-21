package structtag

import (
	"encoding/json"
	"reflect"

	"github.com/Pungyeon/required/pkg/required"
)

// TODO : @pungyeon - This is currently not thread safe. A mutex lock or channel is therefore needed, to ensure no race conditions are met. The reason for this cache implementation, is for general performance. This accounts for a lot of allocations, and since this is static on compilation, we can guarantee that this will never change. Therefore, the cache is a good place to start.
var cache = map[reflect.Type]Tags{}

type Tags struct {
	RequiredInterface  bool
	UnmarshalInterface bool
	Tags               map[string]Tag
}

func (tags Tags) Set(tag Tag) {
	tag.IsSet = true
	tags.Tags[tag.FieldName] = tag
}

func (tags Tags) CheckRequired() error {
	for _, tag := range tags.Tags {
		if tag.Required && !tag.IsSet {
			return requiredErr{
				err:   errRequiredField,
				field: tag.FieldName,
			}
		}
	}
	return nil
}

func (tags Tags) Reset() {
	for _, tag := range tags.Tags {
		tag.IsSet = false
	}
}

func FromValue(vo reflect.Value) (Tags, error) {
	if tags, ok := cache[vo.Type()]; ok {
		tags.Reset()
		return tags, nil
	}

	tags := Tags{Tags: make(map[string]Tag)}
	if vo.CanSet() {
		if vo.CanAddr() {
			_, tags.RequiredInterface = vo.Addr().Interface().(required.Required)
			_, tags.UnmarshalInterface = vo.Addr().Interface().(json.Unmarshaler)
		} else {
			_, tags.RequiredInterface = vo.Interface().(required.Required)
			_, tags.UnmarshalInterface = vo.Interface().(json.Unmarshaler)
		}
	}
	if vo.Kind() != reflect.Struct {
		cache[vo.Type()] = tags
		return tags, nil
	}

	for i := 0; i < vo.NumField(); i++ {
		f := vo.Type().Field(i)
		jsonTag, ok := f.Tag.Lookup("json")
		if !ok {
			tags.Tags[toSnakeCase(f.Name)] = Tag{
				FieldIndex: i,
				FieldName:  f.Name,
			}
		} else {
			tag, err := fromString(jsonTag, i)
			if err != nil {
				return tags, err
			}
			tags.Tags[tag.FieldName] = tag
		}
	}
	cache[vo.Type()] = tags
	return tags, nil
}

var diff uint8 = 'a' - 'A'

func toSnakeCase(s string) string {
	var result string
	for i := 0; i < len(s); i++ {
		if s[i] >= 'A'-1 && s[i] <= 'Z' {
			if i > 0 {
				result += "_"
			}
			result += string(s[i] + diff)
		} else {
			result += string(s[i])
		}
	}
	return result
}
