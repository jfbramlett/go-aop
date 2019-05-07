package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type mdcKey struct {
}

var mdcCtxKey = mdcKey{}

const (
	TIMESTAMP = "timestamp"
	METHOD = "method"
	LEVEL = "level"
	MESSAGE = "msg"
)

func newLogger(ctx context.Context, methodName string, writer LogWriter, config LogConfig) Logger {
	return &logger{ctx: ctx, method: methodName, writer: writer, config: config}
}

type logger struct {
	ctx				context.Context
	method 			string
	writer			LogWriter
	config			LogConfig
}

func (l *logger) IsInfoEnabled() bool {
	return l.isEnabled(InfoLevel)
}

func (l *logger) Info(msg string) {
	if l.IsInfoEnabled() {
		l.logMsg(InfoLevel, msg, nil)
	}
}

func (l *logger) Infof(format string, args ...interface{}) {
	if l.IsInfoEnabled() {
		l.logMsg(InfoLevel, fmt.Sprintf(format, args...), nil)
	}
}

func (l *logger) IsDebugEnabled() bool {
	return l.isEnabled(DebugLevel)
}

func (l *logger) Debug(msg string) {
	if l.IsDebugEnabled() {
		l.logMsg(DebugLevel, msg, nil)
	}
}

func (l *logger) Debugf(format string, args ...interface{}) {
	if l.IsDebugEnabled() {
		l.logMsg(DebugLevel, fmt.Sprintf(format, args...), nil)
	}
}

func (l *logger) IsWarnEnabled() bool {
	return l.isEnabled(WarnLevel)
}

func (l *logger) Warn(msg string) {
	if l.IsWarnEnabled() {
		l.logMsg(WarnLevel, msg, nil)
	}
}

func (l *logger) Warnf(format string, args ...interface{}) {
	if l.IsWarnEnabled() {
		l.logMsg(WarnLevel, fmt.Sprintf(format, args...), nil)
	}
}

func (l *logger) IsErrorEnabled() bool {
	return l.isEnabled(ErrorLevel)
}

func (l *logger) Error(err error, msg string) {
	if l.IsErrorEnabled() {
		l.logMsg(ErrorLevel, msg, err)
	}
}

func (l *logger) Errorf(err error, format string, args ...interface{}) {
	if l.IsErrorEnabled() {
		l.logMsg(ErrorLevel, fmt.Sprintf(format, args...), err)
	}
}

func (l *logger) isEnabled(level string) bool {
	if _, ok := l.config.EnabledLevels[level]; ok {
		return true
	}
	return false
}


func (l *logger) logMsg(level string, msg string, err error) {
	msgJson := make(map[string]interface{})
	msgJson[TIMESTAMP] = time.Now().Format("2006-01-02 15:04:05")
	msgJson[METHOD] = l.method
	msgJson[LEVEL] = level
	msgJson[MESSAGE] = msg
	if err != nil {
		msgJson[ErrorLevel] = err.Error()
	}

	mdc := getMDC(l.ctx)
	for k, v := range mdc {
		msgJson[k] = v
	}

	jsn, err := json.Marshal(msgJson)
	if err == nil {
		l.writer.WriteString(string(jsn))
		return
	}

	l.writer.WriteString(fmt.Sprintf("[%s] [%s] %s", msgJson[TIMESTAMP], level, msg))
}