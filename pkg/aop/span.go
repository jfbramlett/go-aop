package aop

import (
	"context"
	"fmt"
	"github.com/jfbramlett/go-aop/pkg/stackutils"

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
	structName := stackutils.StructNameFromMethod(aop.MethodName)
	methodName := stackutils.MethodNameFromFullPath(aop.MethodName)

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

	span.SetTag(componentKey, component)
	span.SetTag(serviceNameKey, GetServiceName())
	span.SetTag(methodNameKey, stackutils.MethodNameFromFullPath(aop.MethodName))
	span.SetTag(resultKey, result)

	span.Finish()
}
