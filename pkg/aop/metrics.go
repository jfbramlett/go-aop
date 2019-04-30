package aop

import (
	"context"
	"fmt"
	"github.com/namely/go-common/tag"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

const (
	resultTagName  	= "result"
	resultSuccess  	= "success"
	resultFailure  	= "failure"
	component      	= "go-common-timedfunc"
)

var (
	MethodLatencyTimerName = "core_permissions_method_latency"
	MethodLatencyTimerDescription = "Latency metrics for methods invoked within the permissions server"

	ServiceNameTag = tag.MustCreateKey("service_name")
	CallingMethodTag = tag.MustCreateKey("calling_method")
	MethodTag = tag.MustCreateKey("method")
	LatencyTags = []tag.Key {ServiceNameTag, CallingMethodTag, MethodTag}
)


type timerMetricCtxKey struct{}
type spanCtxKey struct{}

var myTimerMetricCtxKey = timerMetricCtxKey{}
var mySpanCtxKey = spanCtxKey{}


var (
	ResultTagKey = tag.MustCreateKey(resultTagName)
)

// NewSpanFuncAspect creats a new Aspect used to wrap something as a new OpenTracing span
func NewSpanFuncAspect(tagKeys []tag.Key) Aspect {
	return &spanAspect{tagKeys: tagKeys}
}

type spanAspect struct {
	tagKeys []tag.Key
}

func (s *spanAspect) Before(ctx context.Context) context.Context {
	aop := AOPFromContext(ctx)
	if aop == nil {
		return ctx
	}

	// establish our span
	span, ctx := opentracing.StartSpanFromContext(ctx, aop.MethodName)
	ext.Component.Set(span, component)

	return context.WithValue(ctx, mySpanCtxKey, span)
}

func (s *spanAspect) After(ctx context.Context, err error) context.Context {
	spanVal := ctx.Value(mySpanCtxKey)
	if spanVal == nil {
		return ctx
	}

	span, valid := spanVal.(opentracing.Span)
	if !valid || span == nil {
		return ctx
	}

	resultCtx := addResultTag(ctx, err)

	tagMap := tag.FromContext(resultCtx)


	values := make([]string, len(s.tagKeys))
	for i, key := range s.tagKeys {
		value, _ := tagMap.Value(key)
		values[i] = value
		span.SetTag(key.Name(), value)
	}


	span.Finish()

	return resultCtx
}

// NewTimedFuncAspect creates a new Aspect that will capture method execution time
func NewTimedFuncAspect(name string, description string, tagKeys []tag.Key) Aspect {

	// Build the set of prometheus labels
	tagKeys = append(tagKeys, ResultTagKey)
	promTags := make([]string, len(tagKeys))
	for i, t := range tagKeys {
		promTags[i] = t.Name()
	}

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

	return &timedFuncAspect{tagKeys: tagKeys, quantiles: quantiles}
}

type timedFuncAspect struct {
	quantiles 	*prometheus.SummaryVec
	tagKeys 	[]tag.Key
}

func (t *timedFuncAspect) Before(ctx context.Context) context.Context {
	aop := AOPFromContext(ctx)
	if aop == nil {
		return ctx
	}

	wrappedContext := context.WithValue(addPrometheusTags(ctx, aop.Service, aop.CallingMethodName, aop.MethodName),
		myTimerMetricCtxKey, time.Now())

	return wrappedContext
}

func (t *timedFuncAspect) After(ctx context.Context, err error) context.Context {
	tStart := ctx.Value(myTimerMetricCtxKey)
	if tStart == nil {
		return ctx
	}

	timerStart, valid := tStart.(time.Time)
	if !valid {
		return ctx
	}

	ms := float64(time.Since(timerStart).Nanoseconds()) / 1e6

	resultCtx := addResultTag(ctx, err)

	tagMap := tag.FromContext(resultCtx)

	values := make([]string, len(t.tagKeys))
	for i, key := range t.tagKeys {
		value, _ := tagMap.Value(key)
		values[i] = value
	}

	// Log the metric
	t.quantiles.WithLabelValues(values...).Observe(ms)

	return resultCtx
}

func addResultTag(ctx context.Context, err error) context.Context {
	result := resultSuccess
	if err != nil {
		result = resultFailure
	}

	// Evaluate the tags
	updatedCtx, ctxError := tag.New(ctx, tag.Insert(ResultTagKey, result))
	if ctxError != nil {
		return ctx
	}

	return updatedCtx
}

func addPrometheusTags(ctx context.Context, source, callingMethod, method string) context.Context {
	taggedCtx, err := tag.AddTagsToContext(ctx, tag.Tag{
		Key:   ServiceNameTag,
		Value: source},
		tag.Tag{
			Key:   CallingMethodTag,
			Value: callingMethod,
		}, tag.Tag{
			Key:   MethodTag,
			Value: method,
		})

	if err != nil {
		return ctx
	}
	return taggedCtx
}
