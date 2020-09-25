package json

import (
	"reflect"
)

func getFieldTags(vo reflect.Value) map[string]int {
	if vo.Kind() != reflect.Struct {
		return map[string]int{}
	}
	tags := make(map[string]int)
	for i := 0; i < vo.NumField(); i++ {
		f := vo.Type().Field(i)
		tag, ok := f.Tag.Lookup("json")
		if !ok {
			tags[toSnakeCase(f.Name)] = i
		} else {
			tags[tag] = i
		}
	}
	return tags
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
