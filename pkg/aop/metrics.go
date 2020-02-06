package aop

import (
	"context"
	"fmt"
	"github.com/jfbramlett/go-aop/pkg/stackutils"
	"github.com/prometheus/client_golang/prometheus"
	"runtime"
	"time"
)

const (
	resultSuccess  		= "success"
	resultFailure  		= "failure"
	component      		= "go-common-timedfunc"
	componentKey		= "component"
	serviceNameKey 		= "service_name"
	callingMethodKey 	= "calling_method"
	methodNameKey 		= "method"
	resultKey 			= "result"
)

type metricCtxKey struct {}
var timerMetricCtxKey = metricCtxKey{}

// NewTimedFuncAdvice creates a new Advice that will capture method execution time
func NewTimedFuncAdvice(name string, description string) Advice {

	// Build the set of prometheus labels
	promTags := []string {serviceNameKey, callingMethodKey, methodNameKey, resultKey}

	// Create the Summary metric
	quantiles := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       fmt.Sprintf("%v_quantiles", name),
			Help:       description,
			Objectives: map[float64]float64{0.5: 0.05, 0.95: 0.005},
		},
		promTags,
	)

	// Register it with the prometheus registry
	err := prometheus.Register(quantiles)
	if err != nil {
		fmt.Println(err)
	}

	return &timedFuncAdvice{quantiles: quantiles}
}

type timedFuncAdvice struct {
	quantiles 	*prometheus.SummaryVec
}

func (t *timedFuncAdvice) Before(ctx context.Context) context.Context {
	aop := AspectFromContext(ctx)
	if aop == nil {
		return ctx
	}

	wrappedContext := context.WithValue(ctx, timerMetricCtxKey, time.Now())

	return wrappedContext
}

func (t *timedFuncAdvice) After(ctx context.Context, err error) {
	aop := AspectFromContext(ctx)
	if aop == nil {
		return
	}

	timerStart, found := t.getStartTime(ctx)
	if !found {
		return
	}

	result := resultSuccess
	if err != nil {
		result = resultFailure
	}

	ms := float64(time.Since(timerStart).Nanoseconds()) / 1e6

	values := []string {GetServiceName(), stackutils.MethodNameFromFullPath(t.getCallingMethod(aop.MethodName)),
		stackutils.MethodNameFromFullPath(aop.MethodName), result}

	// Log the metric
	t.quantiles.WithLabelValues(values...).Observe(ms)
}

func (t *timedFuncAdvice) getStartTime(ctx context.Context) (time.Time, bool) {
	ctxVal := ctx.Value(timerMetricCtxKey)
	if ctxVal != nil {
		if timeStart, ok := ctxVal.(time.Time); ok {
			return timeStart, true
		}
	}
	return time.Time{}, false
}

func (t *timedFuncAdvice) getCallingMethod(toMethod string) string {
	for i := 2;; i++ {
		pc, _, _, ok := runtime.Caller(i)
		details := runtime.FuncForPC(pc)
		if ok && details != nil {
			if details.Name() == toMethod {
				pc, _, _, ok := runtime.Caller(i+1)
				details := runtime.FuncForPC(pc)
				if ok && details != nil {
					return details.Name()
				}
				break
			}
		} else {
			break
		}
	}

	return UnknownMethod
}

