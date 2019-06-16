package aop

import (
	"context"
	"fmt"
	"github.com/jfbramlett/go-aop/pkg/common"
	"github.com/jfbramlett/go-aop/pkg/logging"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
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
	logger := logging.GetLogger(ctx)
	logger.Debug("starting")
	return ctx
}

func (s *loggingAdvice) After(ctx context.Context, err error) {
	logger := logging.GetLogger(ctx)
	if err != nil {
		logger.Debugf("completed with error %s", err)
	} else {
		logger.Debug("completed")
	}
}