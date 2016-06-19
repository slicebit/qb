package qb

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestAggregates(t *testing.T) {
	col := Column("id", Varchar().Size(36))
	assert.Equal(t, Avg(col), Aggregate("AVG", col))
	assert.Equal(t, Count(col), Aggregate("COUNT", col))
	assert.Equal(t, Sum(col), Aggregate("SUM", col))
	assert.Equal(t, Min(col), Aggregate("MIN", col))
	assert.Equal(t, Max(col), Aggregate("MAX", col))
}
