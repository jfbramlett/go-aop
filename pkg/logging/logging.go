package logging

import (
	"context"
	"github.com/sirupsen/logrus"
)

type logKeyStruct struct{}

var logKey = &logKeyStruct{}


// LoggerFromContex get the current logger from the context creating a new one if there is not one in the context
func LoggerFromContext(ctx context.Context) (*logrus.Entry, context.Context) {
	logger := ctx.Value(logKey)
	if logger == nil {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logger = logrus.New()
		ctx = context.WithValue(ctx, logKey, logger)
	}

	return logger.(*logrus.Entry), ctx
}

// ContextWithLogger adds the given logger to the context returning the updated context
func ContextWithLogger(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, logKey, logger)
}

// UpdateInContext updates the logger in the context with the given set of fields
func UpdateInContext(ctx context.Context, fields logrus.Fields) (*logrus.Entry, context.Context) {
	logger, ctx := LoggerFromContext(ctx)
	logger = logger.WithFields(fields)
	return logger, ContextWithLogger(ctx, logger)
}