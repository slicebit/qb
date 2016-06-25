package qb

import (
	"fmt"
	"strings"
)

// Tag is the base abstraction of qb tag
type Tag struct {
	// contains default, null, notnull, unique, primary_key, foreign_key(table.column), check(condition > 0)
	Constraints []string

	// contains type(size) or type parameters
	Type string

	// true if it is ignored
	Ignore bool
}

// ParseTag parses raw qb tag and builds a Tag object
func ParseTag(rawTag string) (Tag, error) {
	rawTag = strings.Replace(rawTag, " ", "", -1)
	rawTag = strings.TrimRight(rawTag, ";")

	tag := Tag{
		Constraints: []string{},
	}

	if rawTag == "" {
		return Tag{}, nil
	}

	tags := strings.Split(rawTag, ";")
	for _, t := range tags {

		tagKeyVal := strings.Split(t, ":")

		if tagKeyVal[0] == "index" {

			tag.Constraints = append(tag.Constraints, t)
			continue

		} else if tagKeyVal[0] == "-" {
			tag.Ignore = true
			return Tag{Ignore: true}, nil
		}

		if len(tagKeyVal) != 2 {
			return Tag{}, fmt.Errorf("Invalid tag key length, tag: %v", tag)
		}

		if tagKeyVal[0] == "type" {
			tag.Type = tagKeyVal[1]
		} else if tagKeyVal[0] == "constraints" || tagKeyVal[0] == "constraint" {
			for _, c := range strings.Split(tagKeyVal[1], ",") {
				if c != "" {
					tag.Constraints = append(tag.Constraints, c)
				}
			}
		} else {
			return Tag{}, fmt.Errorf("Invalid tag key=%s value=%s", tagKeyVal[0], tagKeyVal[1])
		}
	}

	return tag, nil
}

// ParseDBTag parses the "db" tag that can be used in custom column name mapping
func ParseDBTag(rawTag string) string {
	rawTag = strings.Replace(rawTag, " ", "", -1)
	return rawTag
}
