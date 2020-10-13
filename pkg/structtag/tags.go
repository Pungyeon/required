package structtag

import "reflect"

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
