/**
 * Copyright © 2022 Hamed Yousefi <hdyousefi@gmail.com>.
 */

package log

// Logger is an interface for multi-level logging.
type Logger interface {
	Info(msg string, keyValues ...interface{})
	Error(msg string, keyValues ...interface{})
	Warn(msg string, keyValues ...interface{})
	Debug(msg string, keyValues ...interface{})
}
