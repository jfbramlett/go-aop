package aop

import (
	"context"
	"github.com/jfbramlett/go-aop/pkg/logging"
	"github.com/jfbramlett/go-aop/pkg/stackutils"
	"github.com/sirupsen/logrus"
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
	method := stackutils.BasicQualifierFromMethod(ctx.Value(Method).(string))

	logger, newCtx := logging.UpdateInContext(ctx, logrus.Fields{"name": method})
	logger.Debug("starting")
	return newCtx
}

func (s *loggingAdvice) After(ctx context.Context, err error) {
	logger, _ := logging.LoggerFromContext(ctx)
	if err != nil {
		logger.Debugf("completed with error %s", err)
	} else {
		logger.Debug("completed")
	}
}