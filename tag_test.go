package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTag(t *testing.T) {

	tag, _ := ParseTag("type:varchar(255);constraints:default(guest),notnull")
	assert.Equal(t, tag.Type, "varchar(255)")
	assert.Equal(t, tag.Constraints, []string{"default(guest)", "notnull"})

	tagWithoutConstraint, _ := ParseTag("type:varchar(255);constraints:")
	assert.Equal(t, tagWithoutConstraint.Type, "varchar(255)")
	assert.Equal(t, tagWithoutConstraint.Constraints, []string{})

	tagEmpty, _ := ParseTag("     ")
	assert.Equal(t, tagEmpty, &Tag{[]string{}, ""})

	tagInvalidKeyLength, err := ParseTag("type::")
	assert.Equal(t, tagInvalidKeyLength, (*Tag)(nil))
	assert.NotEqual(t, err, nil)
}
