package qb

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	myErr := errors.New("some error")
	stmt := Stmt{}
	qbErr := NewQbError(myErr, &stmt)
	assert.Equal(t, "some error", qbErr.Error())
	assert.Equal(t, &stmt, qbErr.Stmt())
	assert.Equal(t, myErr, qbErr.Orig())
}
