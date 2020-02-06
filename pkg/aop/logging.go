package aop

import (
	"context"
	"github.com/jfbramlett/go-aop/pkg/logging"
)

type loggingCtxKey struct{}
var loggerContextKey = loggingCtxKey{}


// NewLoggingFuncAdvice creats a new Advice used to wrap something as a new OpenTracing span
func NewLoggingFuncAdvice() Advice {
	return &loggingAdvice{}
}

type loggingAdvice struct {
}

func (s *loggingAdvice) Before(ctx context.Context) context.Context {
	method := ctx.Value(Method).(string)
	logger, newCtx := logging.Named(method).Create(ctx)
	logger.Debug("starting")
	return newCtx
}

func (s *loggingAdvice) After(ctx context.Context, err error) {
	method := ctx.Value(Method).(string)
	logger, _ := logging.Named(method).Create(ctx)
	if err != nil {
		logger.Debugf("completed with error %s", err)
	} else {
		logger.Debug("completed")
	}
}