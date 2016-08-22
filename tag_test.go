package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTag(t *testing.T) {
	tag, _ := ParseTag("type:varchar(255);constraints:default(guest),notnull")
	assert.Equal(t, "varchar(255)", tag.Type)
	assert.Equal(t, []string{"default(guest)", "notnull"}, tag.Constraints)

	tagWithoutConstraint, _ := ParseTag("type:varchar(255);constraints:")
	assert.Equal(t, "varchar(255)", tagWithoutConstraint.Type)
	assert.Equal(t, []string{}, tagWithoutConstraint.Constraints)

	tagEmpty, _ := ParseTag("     ")
	assert.Zero(t, tagEmpty)

	tagInvalidKeyLength, err := ParseTag("type::")
	assert.Zero(t, tagInvalidKeyLength)
	assert.NotNil(t, err)
}
