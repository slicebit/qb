package qb

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorMessages(t *testing.T) {
	var tests = []struct {
		code     ErrorCode
		expected string
	}{
		{ErrAny, "Uncategorized error: xxx"},
		{ErrInterface, "Interface error: xxx"},
		{ErrDatabase, "Database error: xxx"},
		{ErrData, "Database data error: xxx"},
		{ErrOperational, "Database operational error: xxx"},
		{ErrIntegrity, "Database integrity error: xxx"},
		{ErrInternal, "Database internal error: xxx"},
		{ErrProgramming, "Database programming error: xxx"},
		{54, "xxx"},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.expected, Error{Code: tt.code, Orig: errors.New("xxx")}.Error())
	}
}

func TestErrorCode(t *testing.T) {
	assert.True(t, ErrInterface.IsInterfaceError())
	assert.False(t, ErrInterface.IsDatabaseError())

	assert.True(t, ErrDatabase.IsDatabaseError())
	assert.False(t, ErrDatabase.IsInterfaceError())

	assert.True(t, ErrData.IsDatabaseError())
	assert.False(t, ErrData.IsInterfaceError())

	assert.True(t, ErrOperational.IsDatabaseError())
	assert.False(t, ErrOperational.IsInterfaceError())

	assert.True(t, ErrIntegrity.IsDatabaseError())
	assert.False(t, ErrIntegrity.IsInterfaceError())

	assert.True(t, ErrInternal.IsDatabaseError())
	assert.False(t, ErrInternal.IsInterfaceError())

	assert.True(t, ErrProgramming.IsDatabaseError())
	assert.False(t, ErrProgramming.IsInterfaceError())
}
