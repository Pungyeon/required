package structtag

import (
	"reflect"
)

var cache = map[reflect.Type]Tags{}

type Tags map[string]Tag

func (tags Tags) Set(tag Tag) {
	tag.IsSet = true
	tags[tag.FieldName] = tag
}

func (tags Tags) CheckRequired() error {
	for _, tag := range tags {
		if tag.Required && !tag.IsSet {
			return requiredErr{
				err:   errRequiredField,
				field: tag.FieldName,
			}
		}
	}
	return nil
}

func FromValue(vo reflect.Value) (Tags, error) {
	if vo.Kind() != reflect.Struct {
		return map[string]Tag{}, nil
	}

	// TODO : @pungyeon - This is currently not thread safe. A mutex lock or channel is therefore needed, to ensure no race conditions are met. The reason for this cache implementation, is for general performance. This accounts for a lot of allocations, and since this is static on compilation, we can guarantee that this will never change. Therefore, the cache is a good place to start.
	// TODO : @pungyeon - On further thought, this might actually break functionality!
	if tags, ok := cache[vo.Type()]; ok {
		return tags, nil
	}

	tags := Tags(make(map[string]Tag))
	for i := 0; i < vo.NumField(); i++ {
		f := vo.Type().Field(i)
		jsonTag, ok := f.Tag.Lookup("json")
		if !ok {
			tags[toSnakeCase(f.Name)] = Tag{
				FieldIndex: i,
				FieldName:  f.Name,
			}
		} else {
			tag, err := fromString(jsonTag, i)
			if err != nil {
				return nil, err
			}
			tags[tag.FieldName] = tag
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
