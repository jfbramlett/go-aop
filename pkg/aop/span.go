package aop

import (
	"context"
	"fmt"

	"github.com/jfbramlett/go-aop/pkg/common"
	"github.com/jfbramlett/go-aop/pkg/tracing"
)

// NewSpanFuncAdvice creates a new Advice used to wrap something as a new OpenTracing span
func NewSpanFuncAdvice() Advice {
	return &spanAdvice{}
}

type spanAdvice struct {
}

func (s *spanAdvice) Before(ctx context.Context) context.Context {
	aop := AspectFromContext(ctx)
	if aop == nil {
		return ctx

	}

	// establish our span
	structName := common.StructNameFromMethod(aop.MethodName)
	methodName := common.MethodNameFromFullPath(aop.MethodName)

	if structName != "" {
		methodName = fmt.Sprintf("%s.%s", structName, methodName)
	}

	_, spanCtx := tracing.StartSpanFromContext(ctx, methodName)

	return spanCtx
}

func (s *spanAdvice) After(ctx context.Context, err error) {
	aop := AspectFromContext(ctx)
	if aop == nil {
		return
	}

	span := tracing.SpanFromContext(ctx)
	if span == nil {
		return
	}

	result := resultSuccess
	if err != nil {
		result = resultFailure
	}

	span.Tag(serviceNameKey, GetServiceName())
	span.Tag(methodNameKey, common.MethodNameFromFullPath(aop.MethodName))
	span.Tag(resultKey, result)

	span.Finish()
}
