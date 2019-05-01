package aop

import (
	"context"
	"errors"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestTimedFuncAdvice(t *testing.T) {
	t.Run("run_success", func(t *testing.T) {
		// given
		callingMethod := "func1"
		method := "TimedMethod1"

		serviceName := "timedFuncSuccess"
		InitAOP(serviceName)

		metricName := "testTimedFuncAdvice"
		RegisterJoinPoint(NewRegexPointcut(".*TimedMethod\\d"), NewTimedFuncAdvice(metricName, "for testing"))

		tStruct := metricsTestSampleStruct{}

		// when
		_, err := tStruct.TimedMethod1(context.Background())

		// then
		assert.Nil(t, err)

		validateMetrics(t, serviceName, []string {callingMethod}, []string{method}, metricName, 1, 0)
	})

	t.Run("run_success_calling_child", func(t *testing.T) {
		// given
		callingMethod := "func2"
		method := "TimedMethod3"
		childMethod := "TimedMethod4"

		serviceName := "timedFuncSuccessWithChild"
		InitAOP(serviceName)

		metricName := "testTimedFuncAdviceWithChild"
		RegisterJoinPoint(NewRegexPointcut(".*TimedMethod\\d"), NewTimedFuncAdvice(metricName, "for testing"))

		tStruct := metricsTestSampleStruct{}

		// when
		_, err := tStruct.TimedMethod3(context.Background())

		// then
		assert.Nil(t, err)

		validateMetrics(t, serviceName, []string {callingMethod, method}, []string {method, childMethod}, metricName, 2, 0)
	})

	t.Run("run_error", func(t *testing.T) {
		// given
		callingMethod := "func3"
		method := "TimedMethod2"

		serviceName := "timedFuncError"
		InitAOP(serviceName)

		metricName := "testTimedFuncAdviceError"
		RegisterJoinPoint(NewRegexPointcut(".*TimedMethod\\d"), NewTimedFuncAdvice(metricName, "for testing"))

		tStruct := metricsTestSampleStruct{}

		// when
		_, err := tStruct.TimedMethod2(context.Background())

		// then
		assert.NotNil(t, err)

		validateMetrics(t, serviceName, []string {callingMethod}, []string {method}, metricName, 0, 1)
	})

}

func TestSpanFuncAdvice(t *testing.T) {
	t.Run("run_success", func(t *testing.T) {
		// given
		mockTracer := &mocktracer.MockTracer{}
		opentracing.SetGlobalTracer(mockTracer)

		expectedOperationName := "SpanMethod1"

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
			methodNameKey: expectedOperationName,
			resultKey: resultSuccess,
		}, startTime, finishTime)
	})

	t.Run("run_error", func(t *testing.T) {
		// given
		mockTracer := &mocktracer.MockTracer{}
		opentracing.SetGlobalTracer(mockTracer)

		expectedOperationName := "SpanMethod2"

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
			methodNameKey: expectedOperationName,
			resultKey: resultFailure,
		}, startTime, finishTime)
	})

	t.Run("run_nested", func(t *testing.T) {
		// given
		mockTracer := &mocktracer.MockTracer{}
		opentracing.SetGlobalTracer(mockTracer)

		expectedOperationName0 := "SpanMethod4"
		expectedOperationName1 := "SpanMethod3"

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
			methodNameKey: expectedOperationName0,
			resultKey: resultSuccess,
		}, finishedSpans[1].StartTime, finishedSpans[1].FinishTime)
		validateSpan(t, finishedSpans[1], expectedOperationName1, map[string]string {"component": component,
			serviceNameKey: serviceName,
			methodNameKey: expectedOperationName1,
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

type metricsTestSampleStruct struct {
	collector		*aspectCollector
}

func (s *metricsTestSampleStruct) TimedMethod1(ctx context.Context) (string, error) {
	var err error
	ctx = Before(ctx, GetMethodName())
	defer func() {After(ctx, err)}()

	return "success", nil
}

func (s *metricsTestSampleStruct) TimedMethod2(ctx context.Context) (string, error) {
	var err error
	ctx = Before(ctx, GetMethodName())
	defer func() {After(ctx, err)}()

	err = errors.New("failed")

	return "", err
}

func (s *metricsTestSampleStruct) TimedMethod3(ctx context.Context) (string, error) {
	var err error
	ctx = Before(ctx, GetMethodName())
	defer func() {After(ctx, err)}()

	return s.TimedMethod4(ctx)
}

func (s *metricsTestSampleStruct) TimedMethod4(ctx context.Context) (string, error) {
	var err error
	ctx = Before(ctx, GetMethodName())
	defer func() {After(ctx, err)}()

	return "success", nil
}

func (s *metricsTestSampleStruct) SpanMethod1(ctx context.Context) (string, error) {
	var err error
	ctx = Before(ctx, GetMethodName())
	defer func() {After(ctx, err)}()

	return "success", nil
}

func (s *metricsTestSampleStruct) SpanMethod2(ctx context.Context) (string, error) {
	var err error
	ctx = Before(ctx, GetMethodName())
	defer func() {After(ctx, err)}()

	err = errors.New("failed")

	return "", err
}

func (s *metricsTestSampleStruct) SpanMethod3(ctx context.Context) (string, error) {
	var err error
	ctx = Before(ctx, GetMethodName())
	defer func() {After(ctx, err)}()

	return s.SpanMethod4(ctx)
}

func (s *metricsTestSampleStruct) SpanMethod4(ctx context.Context) (string, error) {
	var err error
	ctx = Before(ctx, GetMethodName())
	defer func() {After(ctx, err)}()

	return "success", nil
}


func validateMetrics(t *testing.T, serviceName string, callingMethodNames []string, methodNames []string, metricName string, expectSuccess, expectedFailures int) {
	metrics, err := prometheus.DefaultGatherer.Gather()
	require.Nil(t, err)

	latencyFound := false
	for _, metricFamily := range metrics {
		if *metricFamily.Name == metricName + "_quantiles" {
			latencyFound = true
			metricsOfInterest := getMetricsOfInterest(metricFamily, serviceName)
			assert.Equal(t, len(methodNames), len(metricsOfInterest))
			failedCalls := 0
			passedCalls := 0
			for _, metric := range metricsOfInterest {
				assert.True(t, doesLabelMatch(metric, serviceNameKey, serviceName))
				assert.True(t, isLabelInSet(metric, callingMethodKey, callingMethodNames))
				assert.True(t, isLabelInSet(metric, methodNameKey, methodNames))

				if doesLabelMatch(metric, "result", "success") {
					passedCalls++
				} else if doesLabelMatch(metric, "result", "failure") {
					failedCalls++
				}
			}
			assert.Equal(t, expectedFailures, failedCalls)
			assert.Equal(t, expectSuccess, passedCalls)
		}
	}

	assert.True(t, latencyFound)
}

func doesLabelMatch(metric *io_prometheus_client.Metric, labelName string, labelValue string) bool {
	label := getLabel(metric, labelName)
	if label == nil {
		return false
	}
	return labelValue == label.GetValue()
}

func isLabelInSet(metric *io_prometheus_client.Metric, labelName string, labelValues []string) bool {
	label := getLabel(metric, labelName)
	if label == nil {
		return false
	}

	for _, expectedValue := range labelValues {
		if expectedValue == label.GetValue() {
			return true
		}
	}
	return false
}

func getMetricsOfInterest(metricFamily *io_prometheus_client.MetricFamily, serviceName string) []*io_prometheus_client.Metric {
	result := make([]*io_prometheus_client.Metric, 0)

	for _, metric := range metricFamily.Metric {
		label := getLabel(metric, serviceNameKey)
		if label != nil && label.GetValue() == serviceName {
			result = append(result, metric)
		}
	}

	return result
}

func getLabel(metric *io_prometheus_client.Metric, labelName string) *io_prometheus_client.LabelPair {
	for _, label := range metric.Label {
		if label.GetName() == labelName {
			return label
		}
	}
	return nil
}