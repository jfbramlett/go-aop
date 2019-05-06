package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"time"
)

type mdcKey struct {
}

var mdcCtxKey = mdcKey{}

const (
	INFO = "info"
	DEBUG = "debug"
	WARN = "warn"
	ERROR = "error"

	TIMESTAMP = "timestamp"
	METHOD = "method"
	LEVEL = "level"
	MESSAGE = "msg"
)


func AddMDC(ctx context.Context, vals map[string]interface{}) context.Context {
	current := ctx.Value(mdcCtxKey)
	var currentMdc map[string]interface{}
	if current == nil {
		currentMdc = make(map[string]interface{})
	} else {
		currentMdc = current.(map[string]interface{})
	}
	for k, v := range vals {
		currentMdc[k] = v
	}

	return context.WithValue(ctx, mdcCtxKey, currentMdc)
}

func getMDC(ctx context.Context) map[string]interface{} {
	current := ctx.Value(mdcCtxKey)
	if current == nil {
		return make(map[string]interface{})
	} else {
		return current.(map[string]interface{})
	}
}


type Logger interface {
	Info(msg string)
	Infof(fmt string, args ...interface{})
	Debug(msg string)
	Debugf(fmt string, args ...interface{})
	Warn(msg string)
	Warnf(fmt string, args ...interface{})
	Error(err error, msg string)
	Errorf(err error, fmt string, args ...interface{})
}

func GetLogger(ctx context.Context) Logger {
	methodName := getCallingMethodName()
	return &logger{ctx: ctx, method: methodName}
}

type logger struct {
	ctx		context.Context
	method 	string
}

func (l *logger) Info(msg string) {
	l.logMsg(INFO, msg, nil)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.logMsg(INFO, fmt.Sprintf(format, args...), nil)
}

func (l *logger) Debug(msg string) {
	l.logMsg(DEBUG, msg, nil)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.logMsg(DEBUG, fmt.Sprintf(format, args...), nil)
}

func (l *logger) Warn(msg string) {
	l.logMsg(WARN, msg, nil)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.logMsg(WARN, fmt.Sprintf(format, args...), nil)
}

func (l *logger) Error(err error, msg string) {
	l.logMsg(ERROR, msg, err)
}

func (l *logger) Errorf(err error, format string, args ...interface{}) {
	l.logMsg(ERROR, fmt.Sprintf(format, args...), err)
}

func (l *logger) logMsg(level string, msg string, err error) {
	if IsEnabled(level, l.method) {
		msgJson := make(map[string]interface{})
		msgJson[TIMESTAMP] = time.Now().Format("2006-01-02 15:04:05")
		msgJson[METHOD] = l.method
		msgJson[LEVEL] = level
		msgJson[MESSAGE] = msg
		if err != nil {
			msgJson[ERROR] = err.Error()
		}

		mdc := getMDC(l.ctx)
		for k, v := range mdc {
			msgJson[k] = v
		}

		jsn, err := json.Marshal(msgJson)
		if err == nil {
			WriteString(string(jsn))
			return
		}

		WriteString(fmt.Sprintf("[%s] [%s] %s", msgJson[TIMESTAMP], level, msg))
	}
}

func getCallingMethodName() string {
	pc, _, _, ok := runtime.Caller(2)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		return details.Name()
	}

	return "unknown"
}

