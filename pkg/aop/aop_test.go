package aop

import (
	"context"
	"github.com/golang-collections/collections/stack"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const MethodFrame	= "method"
const BeforeFrame	= "before"
const AfterFrame 	= "after"
const LoggingAdvice = "logging"
const CountAdvice	= "count"
const MethodAdvice	= "none"

func TestAOP(t *testing.T) {
	// given
	collector := &aspectCollector{methodCalls: make([]methodCall, 0)}

	InitAOP("testAop")
	RegisterJoinPoint(NewRegexPointcut(".*Method\\d$"), &loggingAspect{collector: collector})
	RegisterJoinPoint(NewRegexPointcut(".*Method1$"), &countingAspect{collector: collector})

	expected := []methodCall{{BeforeFrame, "github.com/jfbramlett/go-aop/pkg/aop.(*sampleStruct).Method1", LoggingAdvice},
		{BeforeFrame, "github.com/jfbramlett/go-aop/pkg/aop.(*sampleStruct).Method1", CountAdvice},
		{MethodFrame, "Method1", MethodAdvice},
		{AfterFrame, "github.com/jfbramlett/go-aop/pkg/aop.(*sampleStruct).Method1", CountAdvice},
		{AfterFrame, "github.com/jfbramlett/go-aop/pkg/aop.(*sampleStruct).Method1", LoggingAdvice},
		{BeforeFrame, "github.com/jfbramlett/go-aop/pkg/aop.(*sampleStruct).Method2", LoggingAdvice},
		{MethodFrame, "Method2", MethodAdvice},
		{AfterFrame, "github.com/jfbramlett/go-aop/pkg/aop.(*sampleStruct).Method2", LoggingAdvice},
		{BeforeFrame, "github.com/jfbramlett/go-aop/pkg/aop.(*sampleStruct).Method3", LoggingAdvice},
		{BeforeFrame, "github.com/jfbramlett/go-aop/pkg/aop.(*sampleStruct).privateMethod1", LoggingAdvice},
		{BeforeFrame, "github.com/jfbramlett/go-aop/pkg/aop.(*sampleStruct).privateMethod1", CountAdvice},
		{MethodFrame, "privateMethod1", MethodAdvice},
		{AfterFrame, "github.com/jfbramlett/go-aop/pkg/aop.(*sampleStruct).privateMethod1", CountAdvice},
		{AfterFrame, "github.com/jfbramlett/go-aop/pkg/aop.(*sampleStruct).privateMethod1", LoggingAdvice},
		{AfterFrame, "github.com/jfbramlett/go-aop/pkg/aop.(*sampleStruct).Method3", LoggingAdvice},
		{MethodFrame, "Special", MethodAdvice},
	}

	// when
	st := sampleStruct{collector: collector}
	st.Method1("arg1", 1)
	st.Method2("arg1", 1)
	st.Method3("arg1", 1)
	st.Special("arg1", 1)

	// then
	assert.Equal(t, expected, collector.methodCalls)
}

func TestGetServiceName(t *testing.T) {
	// given
	expectedServiceName := uuid.New().String()
	InitAOP(expectedServiceName)

	// when
	serviceName := GetServiceName()

	// then
	assert.Equal(t, expectedServiceName, serviceName)
}


func TestGetMethodName(t *testing.T) {
	// given
	expectedMethodName := "github.com/jfbramlett/go-aop/pkg/aop.TestGetMethodName"

	// when
	methodName := GetMethodName()

	// then
	assert.Equal(t, expectedMethodName, methodName)
}

func TestAspectFromContext(t *testing.T) {
	t.Run("aspect_exist", func(t *testing.T) {
		// given
		expectedAspect := &Aspect{}
		mgr := aspectMgr{}
		ctx := mgr.contextWithAspect(context.Background(), expectedAspect)

		// when
		aspect := AspectFromContext(ctx)

		// then
		assert.Equal(t, expectedAspect, aspect)
	})

	t.Run("aspect_exist_multi", func(t *testing.T) {
		// given
		initialAspect := &Aspect{}
		expectedAspect := &Aspect{}
		mgr := aspectMgr{}
		ctx := mgr.contextWithAspect(context.Background(), initialAspect)
		ctx = mgr.contextWithAspect(ctx, expectedAspect)

		// when
		poppedAspect := AspectFromContext(ctx)
		newCtx := mgr.removeAspectFromContext(ctx)
		aspect := AspectFromContext(newCtx)

		// then
		assert.Equal(t, expectedAspect, poppedAspect)
		assert.Equal(t, initialAspect, aspect)
	})

	t.Run("no_aspect", func(t *testing.T) {
		// given
		ctx := context.Background()

		// when
		aspect := AspectFromContext(ctx)

		// then
		assert.Nil(t, aspect)
	})

}

func TestMethodNameFromFullPath(t *testing.T) {
	t.Run("test_full_name", func(t *testing.T) {
		// given
		expectedMethodName := "MyMethod"
		fullPath := "github.com/jfbramlett/go-aop/pkg/metrics." + expectedMethodName

		// when
		methodName := MethodNameFromFullPath(fullPath)

		// then
		assert.Equal(t, expectedMethodName, methodName)
	})

	t.Run("test_malformed_name", func(t *testing.T) {
		// given
		expectedMethodName := "MyMethod"

		// when
		methodName := MethodNameFromFullPath(expectedMethodName)

		// then
		assert.Equal(t, expectedMethodName, methodName)
	})

}

func TestPushToContext(t *testing.T) {
	t.Run("add_to_empty", func(t *testing.T) {
		// given
		key := "ctxKey"
		val := "some value"

		// when
		ctx := PushToContext(context.Background(), key, val)

		// then
		ctxVal := ctx.Value(key)
		require.NotNil(t, ctxVal)

		stackVal, valid := ctxVal.(*stack.Stack)
		require.True(t, valid)
		assert.Equal(t, 1, stackVal.Len())
		assert.Equal(t, val, stackVal.Pop())
	})

	t.Run("value_already_there", func(t *testing.T) {
		// given
		stackVal := stack.New()
		key := "ctxKey"
		val1 := "some value"
		val2 := "some other value"
		stackVal.Push(val1)
		ctx := context.Background()
		ctx = context.WithValue(ctx, key, stackVal)

		// when
		ctx = PushToContext(ctx, key, val2)

		// then
		ctxVal := ctx.Value(key)
		require.NotNil(t, ctxVal)

		stackVal, valid := ctxVal.(*stack.Stack)
		require.True(t, valid)
		assert.Equal(t, 2, stackVal.Len())
		assert.Equal(t, val2, stackVal.Pop())
		assert.Equal(t, val1, stackVal.Pop())
	})
}

func TestPopFromContext(t *testing.T) {
	t.Run("pop_empty", func(t *testing.T) {
		// given
		key := "ctxKey"

		// when
		ctx, val := PopFromContext(context.Background(), key)

		// then
		assert.NotNil(t, ctx)
		assert.Nil(t, val)
	})

	t.Run("pop_value", func(t *testing.T) {
		// given
		stackVal := stack.New()
		key := "ctxKey"
		val1 := "some value"
		val2 := "some other value"
		stackVal.Push(val1)
		stackVal.Push(val2)
		ctx := context.Background()
		ctx = context.WithValue(ctx, key, stackVal)

		// when
		ctx, val := PopFromContext(ctx, key)

		// then
		assert.NotNil(t, ctx)
		assert.Equal(t, val2, val)
	})

}


type sampleStruct struct {
	collector		*aspectCollector
}

func (s *sampleStruct) Method1(arg1 string, arg2 int) (string, error) {
	var err error
	ctx := Before(context.Background(), GetMethodName())
	defer func() {After(ctx, err)}()

	s.collector.Collect(MethodFrame, "Method1", MethodAdvice)

	return "success", nil
}

func (s *sampleStruct) Method2(arg1 string, arg2 int) (string, error) {
	var err error
	ctx := Before(context.Background(), GetMethodName())
	defer func() {After(ctx, err)}()

	s.collector.Collect(MethodFrame, "Method2", MethodAdvice)

	return "success", nil
}

func (s *sampleStruct) Method3(arg1 string, arg2 int) (string, error) {
	var err error
	ctx := Before(context.Background(), GetMethodName())
	defer func() {After(ctx, err)}()

	s.privateMethod1(arg1, arg2)

	return "success", nil
}

func (s *sampleStruct) Special(arg1 string, arg2 int) (string, error) {
	var err error
	ctx := Before(context.Background(), GetMethodName())
	defer func() {After(ctx, err)}()

	s.collector.Collect(MethodFrame, "Special", MethodAdvice)

	return "success", nil
}

func (s *sampleStruct) privateMethod1(arg1 string, arg2 int) (string, error) {
	var err error
	ctx := Before(context.Background(), GetMethodName())
	defer func() {After(ctx, err)}()

	s.collector.Collect(MethodFrame, "privateMethod1", MethodAdvice)

	return "success", nil
}


type loggingAspect struct {
	collector		*aspectCollector
}

func (l *loggingAspect) Before(ctx context.Context) context.Context {
	definition := AspectFromContext(ctx)
	l.collector.Collect(BeforeFrame, definition.MethodName, LoggingAdvice)
	return ctx
}

func (l *loggingAspect) After(ctx context.Context, err error) context.Context {
	definition := AspectFromContext(ctx)
	l.collector.Collect(AfterFrame, definition.MethodName, LoggingAdvice)
	return ctx
}

type countingAspect struct {
	collector		*aspectCollector
}

func (c *countingAspect) Before(ctx context.Context) context.Context {
	definition := AspectFromContext(ctx)
	c.collector.Collect(BeforeFrame, definition.MethodName, CountAdvice)
	return ctx
}

func (c *countingAspect) After(ctx context.Context, err error) context.Context {
	definition := AspectFromContext(ctx)
	c.collector.Collect(AfterFrame, definition.MethodName, CountAdvice)
	return ctx
}


type methodCall struct {
	frame					string
	methodName				string
	op						string
}

type aspectCollector struct {
	methodCalls		[]methodCall
}

func (a *aspectCollector) Collect(frame, methodName, op string) *aspectCollector {
	a.methodCalls = append(a.methodCalls, methodCall{frame, methodName, op})
	return a
}



