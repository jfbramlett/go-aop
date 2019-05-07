package logging

const (
	InfoLevel  = "info"
	DebugLevel = "debug"
	WarnLevel  = "warn"
	ErrorLevel = "error"
)

type Logger interface {
	IsInfoEnabled() bool
	Info(msg string)
	Infof(fmt string, args ...interface{})

	IsDebugEnabled() bool
	Debug(msg string)
	Debugf(fmt string, args ...interface{})

	IsWarnEnabled() bool
	Warn(msg string)
	Warnf(fmt string, args ...interface{})

	IsErrorEnabled() bool
	Error(err error, msg string)
	Errorf(err error, fmt string, args ...interface{})
}

type LogConfig struct {
	EnabledLevels		map[string]bool
}

var DefaultLogConfig = LogConfig{EnabledLevels: map[string]bool {InfoLevel: true, DebugLevel: true, WarnLevel: true, ErrorLevel: true}}
