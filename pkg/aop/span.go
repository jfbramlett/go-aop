package aop

import (
    "context"
    "fmt"
    "github.com/jfbramlett/go-aop/pkg/common"
    "github.com/opentracing/opentracing-go"
    "github.com/opentracing/opentracing-go/ext"
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

    span, ctx := opentracing.StartSpanFromContext(ctx, methodName)
    ext.Component.Set(span, component)

    return ctx
}

func (s *spanAdvice) After(ctx context.Context, err error) {
    aop := AspectFromContext(ctx)
    if aop == nil {
        return
    }

    span := opentracing.SpanFromContext(ctx)
    if span == nil {
        return
    }

    result := resultSuccess
    if err != nil {
        result = resultFailure
    }

    span.SetTag(serviceNameKey, GetServiceName())
    span.SetTag(methodNameKey, common.MethodNameFromFullPath(aop.MethodName))
    span.SetTag(resultKey, result)

    span.Finish()
}