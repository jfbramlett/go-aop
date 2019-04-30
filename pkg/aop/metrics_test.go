package aop

import (
	"context"
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)


func TestTimedFuncAdvice(t *testing.T) {
	t.Run("run_success", func(t *testing.T) {
		// given
		callingMethod := "func1"
		method := "TimedMethod1"

		serviceName := "timedFuncSuccess"
		InitAOP(serviceName)

		metricName := "testTimedFuncAdvice"
		RegisterJoinPoint(".*TimedMethod\\d", NewTimedFuncAdvice(metricName, "for testing"))

		tStruct := metricsTestSampleStruct{}

		// when
		_, err := tStruct.TimedMethod1("arg1", 1)

		// then
		assert.Nil(t, err)

		validateMetrics(t, serviceName, callingMethod, method, metricName, 1, 0)
	})

	t.Run("run_error", func(t *testing.T) {
		// given
		callingMethod := "func2"
		method := "TimedMethod2"

		serviceName := "timedFuncError"
		InitAOP(serviceName)

		metricName := "testTimedFuncAdviceError"
		RegisterJoinPoint(".*TimedMethod\\d", NewTimedFuncAdvice(metricName, "for testing"))

		tStruct := metricsTestSampleStruct{}

		// when
		_, err := tStruct.TimedMethod2("arg1", 1)

		// then
		assert.NotNil(t, err)

		validateMetrics(t, serviceName, callingMethod, method, metricName, 0, 1)
	})

}

func TestSpanFuncAdvice(t *testing.T) {
	t.Run("run_success", func(t *testing.T) {
		// given
		serviceName := "spanFuncSuccess"
		InitAOP(serviceName)

		RegisterJoinPoint(".*SpanMethod\\d", NewSpanFuncAdvice())

		tStruct := metricsTestSampleStruct{}

		// when
		_, err := tStruct.SpanMethod1("arg1", 1)

		// then
		assert.Nil(t, err)
	})

	t.Run("run_error", func(t *testing.T) {
		// given
		serviceName := "spanFuncError"
		InitAOP(serviceName)

		RegisterJoinPoint(".*SpanMethod\\d", NewSpanFuncAdvice())

		tStruct := metricsTestSampleStruct{}

		// when
		_, err := tStruct.SpanMethod2("arg1", 1)

		// then
		assert.NotNil(t, err)
	})

}

type metricsTestSampleStruct struct {
	collector		*aspectCollector
}

func (s *metricsTestSampleStruct) TimedMethod1(arg1 string, arg2 int) (string, error) {
	var err error
	ctx := Before(context.Background())
	defer func() {After(ctx, err)}()

	return "success", nil
}

func (s *metricsTestSampleStruct) TimedMethod2(arg1 string, arg2 int) (string, error) {
	var err error
	ctx := Before(context.Background())
	defer func() {After(ctx, err)}()

	err = errors.New("failed")

	return "", err
}

func (s *metricsTestSampleStruct) SpanMethod1(arg1 string, arg2 int) (string, error) {
	var err error
	ctx := Before(context.Background())
	defer func() {After(ctx, err)}()

	return "success", nil
}

func (s *metricsTestSampleStruct) SpanMethod2(arg1 string, arg2 int) (string, error) {
	var err error
	ctx := Before(context.Background())
	defer func() {After(ctx, err)}()

	err = errors.New("failed")

	return "", err
}


func validateMetrics(t *testing.T, serviceName, callingMethodName, methodName, metricName string, expectSuccess, expectedFailures int) {
	metrics, err := prometheus.DefaultGatherer.Gather()
	require.Nil(t, err)

	latencyFound := false
	for _, metricFamily := range metrics {
		if *metricFamily.Name == metricName + "_quantiles" {
			latencyFound = true
			metrics := getMetricsOfInterest(metricFamily, serviceName)
			assert.Equal(t, 1, len(metrics))
			failedCalls := 0
			passedCalls := 0
			for _, metric := range metrics {
				assert.True(t, doesLabelMatch(metric, ServiceNameTag.Name(), serviceName))
				assert.True(t, doesLabelMatch(metric, CallingMethodTag.Name(), callingMethodName))
				assert.True(t, doesLabelMatch(metric, MethodTag.Name(), methodName))

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

func getMetricsOfInterest(metricFamily *io_prometheus_client.MetricFamily, serviceName string) []*io_prometheus_client.Metric {
	result := make([]*io_prometheus_client.Metric, 0)

	for _, metric := range metricFamily.Metric {
		label := getLabel(metric, ServiceNameTag.Name())
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