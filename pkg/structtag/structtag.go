package structtag

import (
	"fmt"
)

type Tag struct {
	FieldIndex  int
	FieldName   string
	Required    bool
	OmitIfEmpty bool
	IsSet       bool
}

func (t *Tag) addValue(value string) error {
	if t.FieldName == "" {
		t.FieldName = value
		return nil
	}
	switch value {
	case "required":
		t.Required = true
	case "omitifempty":
		t.OmitIfEmpty = true
	default:
		return fmt.Errorf("illegal tag value: `%s`", value)
	}
	return nil
}

func fromString(input string, index int) (Tag, error) {
	tag := Tag{
		FieldIndex: index,
	}
	var previous int
	current := indexOfNextTag(input, 0)
	for current < len(input) {
		if input[current] == ',' {
			if err := tag.addValue(input[previous:current]); err != nil {
				return tag, err
			}
			current = indexOfNextTag(input, current)
			previous = current
		}
		current++
	}
	if previous != current {
		return tag, tag.addValue(input[previous:current])
	}
	return tag, nil
}

func indexOfNextTag(input string, current int) int {
	for input[current] == ' ' ||
		input[current] == '\n' ||
		input[current] == '\t' ||
		input[current] == ',' ||
		input[current] == '\r' {
		current++
	}
	return current
}
