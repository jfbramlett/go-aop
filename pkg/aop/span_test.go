package aop

import (
    "context"
    "github.com/opentracing/opentracing-go"
    "github.com/opentracing/opentracing-go/mocktracer"
    "github.com/stretchr/testify/assert"
    "testing"
    "time"
)

func TestSpanFuncAdvice(t *testing.T) {
    t.Run("run_success", func(t *testing.T) {
        // given
        mockTracer := &mocktracer.MockTracer{}
        opentracing.SetGlobalTracer(mockTracer)

        expectedOperationName := "metricsTestSampleStruct.SpanMethod1"
        expectedMethodName := "SpanMethod1"

        serviceName := "spanFuncSuccess"
        InitAOP(serviceName)

        RegisterJoinPoint(NewRegexPointcut(".*SpanMethod\\d"), NewSpanFuncAdvice())

        tStruct := metricsTestSampleStruct{}

        // when
        startTime := time.Now()
        _, err := tStruct.SpanMethod1(context.Background())
        finishTime := time.Now()

        // then
        assert.Nil(t, err)

        finishedSpans := mockTracer.FinishedSpans()
        assert.Equal(t, 1, len(finishedSpans))
        validateSpan(t, finishedSpans[0], expectedOperationName, map[string]string {"component": component,
            serviceNameKey: serviceName,
            methodNameKey: expectedMethodName,
            resultKey: resultSuccess,
        }, startTime, finishTime)
    })

    t.Run("run_error", func(t *testing.T) {
        // given
        mockTracer := &mocktracer.MockTracer{}
        opentracing.SetGlobalTracer(mockTracer)

        expectedOperationName := "metricsTestSampleStruct.SpanMethod2"
        expectedMethodName := "SpanMethod2"

        serviceName := "spanFuncError"
        InitAOP(serviceName)

        RegisterJoinPoint(NewRegexPointcut(".*SpanMethod\\d"), NewSpanFuncAdvice())

        tStruct := metricsTestSampleStruct{}

        // when
        startTime := time.Now()
        _, err := tStruct.SpanMethod2(context.Background())
        finishTime := time.Now()

        // then
        assert.NotNil(t, err)

        finishedSpans := mockTracer.FinishedSpans()
        assert.Equal(t, 1, len(finishedSpans))
        validateSpan(t, finishedSpans[0], expectedOperationName, map[string]string {"component": component,
            serviceNameKey: serviceName,
            methodNameKey: expectedMethodName,
            resultKey: resultFailure,
        }, startTime, finishTime)
    })

    t.Run("run_nested", func(t *testing.T) {
        // given
        mockTracer := &mocktracer.MockTracer{}
        opentracing.SetGlobalTracer(mockTracer)

        expectedOperationName0 := "metricsTestSampleStruct.SpanMethod4"
        expectedMethodName0 := "SpanMethod4"
        expectedOperationName1 := "metricsTestSampleStruct.SpanMethod3"
        expectedMethodName1 := "SpanMethod3"

        serviceName := "spanFuncError"
        InitAOP(serviceName)

        RegisterJoinPoint(NewRegexPointcut(".*SpanMethod\\d"), NewSpanFuncAdvice())

        tStruct := metricsTestSampleStruct{}

        // when
        startTime := time.Now()
        _, err := tStruct.SpanMethod3(context.Background())
        finishTime := time.Now()

        // then
        assert.Nil(t, err)

        finishedSpans := mockTracer.FinishedSpans()
        assert.Equal(t, 2, len(finishedSpans))
        validateSpan(t, finishedSpans[0], expectedOperationName0, map[string]string {"component": component,
            serviceNameKey: serviceName,
            methodNameKey: expectedMethodName0,
            resultKey: resultSuccess,
        }, finishedSpans[1].StartTime, finishedSpans[1].FinishTime)
        validateSpan(t, finishedSpans[1], expectedOperationName1, map[string]string {"component": component,
            serviceNameKey: serviceName,
            methodNameKey: expectedMethodName1,
            resultKey: resultSuccess,
        }, startTime, finishTime)

    })
}

func validateSpan(t *testing.T, span *mocktracer.MockSpan, expectedOperationName string, tags map[string]string, timeStartAfter time.Time, timeFinishBefore time.Time) {
    assert.Equal(t, expectedOperationName, span.OperationName)
    assert.True(t, timeStartAfter.Before(span.StartTime))
    assert.True(t, timeFinishBefore.After(span.FinishTime))
    for k, v := range tags {
        assert.Equal(t, v, span.Tag(k))
    }
}

