package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAggregates(t *testing.T) {
	col := Column("id", Varchar().Size(36))
	assert.Equal(t, Aggregate("AVG", col), Avg(col))
	assert.Equal(t, Aggregate("COUNT", col), Count(col))
	assert.Equal(t, Aggregate("SUM", col), Sum(col))
	assert.Equal(t, Aggregate("MIN", col), Min(col))
	assert.Equal(t, Aggregate("MAX", col), Max(col))
}
