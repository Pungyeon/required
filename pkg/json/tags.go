package json

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	errRequiredField = errors.New("required field missing")
)

func IsRequiredErr(err error) bool {
	_, ok := err.(requiredErr)
	return ok
}

type requiredErr struct {
	err   error
	field string
}

func (err requiredErr) Error() string {
	return fmt.Sprintf("%v: %s", err.err, err.field)
}

var _ error = &requiredErr{}

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

type Tag struct {
	FieldIndex  int
	FieldName   string
	Required    bool
	OmitIfEmpty bool
	IsSet       bool
}

func (t *Tag) AddTagValue(value string) error {
	switch value {
	case "required":
		t.Required = true
	case "omitifempty":
		t.OmitIfEmpty = true
	default:
		return fmt.Errorf("illegal tag value: %s", value)
	}
	return nil
}

func getFieldTags(vo reflect.Value) (Tags, error) {
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
			tagValues := strings.Split(jsonTag, ",")
			tag := Tag{
				FieldIndex: i,
				FieldName:  tagValues[0],
			}
			for i := 1; i < len(tagValues); i++ {
				if err := tag.AddTagValue(tagValues[i]); err != nil {
					return nil, err
				}
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
