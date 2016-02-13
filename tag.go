package qbit

import (
	"strings"
	//	"fmt"
	"errors"
	"fmt"
)

// Tag is the base abstraction of qbit tag
type Tag struct {

	// contains default, null, notnull, unique, primary_key, foreign_key(table.column), check(condition > 0)
	Constraints []string

	// contains type(size) or type parameters
	Type string
}

// ParseTag parses raw qbit tag and builds a Tag object
func ParseTag(rawTag string) (*Tag, error) {

	tag := &Tag{
		Constraints: []string{},
	}

	if strings.Trim(rawTag, " ") == "" {
		return tag, nil
	}

	tags := strings.Split(rawTag, ";")
	for _, t := range tags {
		tagKey := strings.Split(t, ":")
		if len(tagKey) == 2 {
			if tagKey[0] == "type" {
				tag.Type = tagKey[1]
			} else if tagKey[0] == "constraints" || tagKey[0] == "constraint" {
				tag.Constraints = strings.Split(tagKey[1], ",")
			} else {
				return nil, errors.New("Invalid keyword in struct tag")
			}
		} else {
			return nil, fmt.Errorf("Invalid tag key length, tag: %v", tag)
		}
	}

	return tag, nil
}
