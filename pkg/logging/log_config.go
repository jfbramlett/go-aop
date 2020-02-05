package logging

import (
	"context"
	"github.com/jfbramlett/go-aop/pkg/common"
	"os"
)

// InitLogging initializes our logging providing a factory function for obtaining new loggers
func InitLogging() {
	globalLogConfig = &logConfig{defaultWriter: initChannelLogWriter(os.Stdout)}
}

// GetLogger is our factory function for obtaining a new logger for the given context. The log config is based on the
// calling method name
func GetLogger(ctx context.Context) Logger {
	callingMethod := common.GetCallingMethodName()
	return globalLogConfig.GetLogger(ctx, common.BasicQualifierFromMethod(callingMethod))
}

// GetLoggerFor is our factory function for obtaining a new logger for the given context with a config based on the given
// free form name
func GetLoggerFor(ctx context.Context, forName string) Logger {
	return globalLogConfig.GetLogger(ctx, forName)
}


type logConfig struct {
	defaultWriter		LogWriter
}

var globalLogConfig *logConfig

func (l *logConfig) GetLogger(ctx context.Context, methodName string) Logger {
	return newLogger(ctx, methodName, l.getWriter(methodName), DefaultLogConfig)
}

func (l *logConfig) getWriter(methodName string) LogWriter {
	return l.defaultWriter
}

