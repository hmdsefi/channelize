/**
 * Copyright © 2022 Hamed Yousefi <hdyousefi@gmail.com>.
 */

package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testCode  = "test_code"
	testError = "test error"
	testField = "test"
	empty     = ""
)

func TestResult_IsValid(t *testing.T) {
	out := new(Result)
	assert.True(t, out.IsValid())

	out.AddFieldError(testField, testError)
	assert.False(t, out.IsValid())

	out = NewResult(testCode, testError)
	assert.False(t, out.IsValid())

	out.AddFieldError(testField, testError)
	assert.False(t, out.IsValid())

	out = NewResult(testCode, empty)
	assert.False(t, out.IsValid())

	out.AddFieldError(testField, testError)
	assert.False(t, out.IsValid())

	out = NewResult(empty, testError)
	assert.False(t, out.IsValid())

	out.AddFieldError(testField, testError)
	assert.False(t, out.IsValid())
}
