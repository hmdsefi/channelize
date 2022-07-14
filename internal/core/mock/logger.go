/**
 * Copyright © 2022 Hamed Yousefi <hdyousefi@gmail.com>.
 */

package mock

import (
	"github.com/hmdsefi/channelize/internal/common/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Logger struct {
	t        *testing.T
	expected []interface{}
}

func NewLogger(t *testing.T) *Logger {
	return &Logger{t: t}
}

func (l *Logger) Expected(expected ...interface{}) {
	l.expected = expected
}

func (l Logger) Info(msg string, keyValues ...interface{}) {
}

func (l Logger) Error(msg string, keyValues ...interface{}) {
	actualArgs := append([]interface{}{log.Error, msg}, keyValues...)
	if len(l.expected) != len(actualArgs) {
		l.t.Fatal("not enough expected argument")
	}
	for i, expected := range l.expected {
		assert.Equal(l.t, expected, actualArgs[i])
	}
}

func (l Logger) Warn(msg string, keyValues ...interface{}) {
}

func (l Logger) Debug(msg string, keyValues ...interface{}) {
}
