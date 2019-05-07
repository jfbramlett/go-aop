package logging

import "context"

func newNoopLogger(ctx context.Context, methodName string) Logger {
	return &noopLogger{}
}

type noopLogger struct {
}

func (l *noopLogger) IsInfoEnabled() bool {
	return false
}

func (l *noopLogger) Info(msg string) {
}

func (l *noopLogger) Infof(format string, args ...interface{}) {
}

func (l *noopLogger) IsDebugEnabled() bool {
	return false
}

func (l *noopLogger) Debug(msg string) {
}

func (l *noopLogger) Debugf(format string, args ...interface{}) {
}

func (l *noopLogger) IsWarnEnabled() bool {
	return false
}

func (l *noopLogger) Warn(msg string) {
}

func (l *noopLogger) Warnf(format string, args ...interface{}) {
}

func (l *noopLogger) IsErrorEnabled() bool {
	return false
}

func (l *noopLogger) Error(err error, msg string) {
}

func (l *noopLogger) Errorf(err error, format string, args ...interface{}) {
}

