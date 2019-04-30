package aop

import (
	"context"
	"github.com/google/uuid"
	"github.com/namely/go-common/tag"
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
	RegisterJoinPoint(".*Method\\d$", &loggingAspect{collector: collector})
	RegisterJoinPoint(".*Method1$", &countingAspect{collector: collector})

	expected := []methodCall{{BeforeFrame, "Method1", "TestAOP", LoggingAdvice},
		{BeforeFrame, "Method1", "TestAOP", CountAdvice},
		{MethodFrame, "Method1", "", MethodAdvice},
		{AfterFrame, "Method1", "TestAOP", CountAdvice},
		{AfterFrame, "Method1", "TestAOP", LoggingAdvice},
		{BeforeFrame, "Method2", "TestAOP", LoggingAdvice},
		{MethodFrame, "Method2", "", MethodAdvice},
		{AfterFrame, "Method2", "TestAOP", LoggingAdvice},
		{BeforeFrame, "Method3", "TestAOP", LoggingAdvice},
		{BeforeFrame, "privateMethod1", "Method3", LoggingAdvice},
		{BeforeFrame, "privateMethod1", "Method3", CountAdvice},
		{MethodFrame, "privateMethod1", "", MethodAdvice},
		{AfterFrame, "privateMethod1", "Method3", CountAdvice},
		{AfterFrame, "privateMethod1", "Method3", LoggingAdvice},
		{AfterFrame, "Method3", "TestAOP", LoggingAdvice},
		{MethodFrame, "Special", "", MethodAdvice},
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


func TestGetMethod(t *testing.T) {
	// given
	expectedMethodName := "github.com/jfbramlett/go-aop/pkg/aop.TestGetMethod"

	// when
	methodName := level1Func()

	// then
	assert.Equal(t, expectedMethodName, methodName)
}

func TestGetCallingMethod(t *testing.T) {
	// given:
	expectedCallingMethodName := "TestGetCallingMethod"

	// when
	callingMethodName := delegateCall()

	// then
	assert.Equal(t, expectedCallingMethodName, callingMethodName)
}

func TestAspectFromContext(t *testing.T) {
	t.Run("aspect_exist", func(t *testing.T) {
		// given
		expectedAspect := &Aspect{}
		ctx := contextWithAspect(context.Background(), expectedAspect)

		// when
		aspect := AspectFromContext(ctx)

		// then
		assert.Equal(t, expectedAspect, aspect)
	})

	t.Run("aspect_exist_multi", func(t *testing.T) {
		// given
		initialAspect := &Aspect{}
		expectedAspect := &Aspect{}
		ctx := contextWithAspect(context.Background(), initialAspect)
		ctx = contextWithAspect(ctx, expectedAspect)

		// when
		poppedAspect := AspectFromContext(ctx)
		newCtx := removeAspectFromContext(ctx)
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

func level1Func() string {
	return getMethod()
}

func delegateCall() string {
	return level2Func()
}

func level2Func() string {
	return getCallingMethod()
}

func TestAddPrometheusTags(t *testing.T) {
	// given
	origCtx := context.Background()
	source := "my source"
	callingMethod := "calling method"
	method := "the method"

	// when
	updatedCtx := addPrometheusTags(origCtx, source, callingMethod, method)

	// then
	assert.NotEqual(t, origCtx, updatedCtx)

	tags := tag.FromContext(updatedCtx)

	assert.NotNil(t, tags)
	value, found := tags.Value(ServiceNameTag)
	require.True(t, found)
	assert.Equal(t, source, value)

	value, found = tags.Value(CallingMethodTag)
	require.True(t, found)
	assert.Equal(t, callingMethod, value)

	value, found = tags.Value(MethodTag)
	require.True(t, found)
	assert.Equal(t, method, value)
}

func TestMethodNameFromFullPath(t *testing.T) {
	t.Run("test_full_name", func(t *testing.T) {
		// given
		expectedMethodName := "MyMethod"
		fullPath := "github.com/namely/permissions/pkg/metrics." + expectedMethodName

		// when
		methodName := methodNameFromFullPath(fullPath)

		// then
		assert.Equal(t, expectedMethodName, methodName)
	})

	t.Run("test_malformed_name", func(t *testing.T) {
		// given
		expectedMethodName := "MyMethod"

		// when
		methodName := methodNameFromFullPath(expectedMethodName)

		// then
		assert.Equal(t, expectedMethodName, methodName)
	})

}

type sampleStruct struct {
	collector		*aspectCollector
}

func (s *sampleStruct) Method1(arg1 string, arg2 int) (string, error) {
	var err error
	ctx := Before(context.Background())
	defer func() {After(ctx, err)}()

	s.collector.Collect(MethodFrame, "Method1", "", MethodAdvice)

	return "success", nil
}

func (s *sampleStruct) Method2(arg1 string, arg2 int) (string, error) {
	var err error
	ctx := Before(context.Background())
	defer func() {After(ctx, err)}()

	s.collector.Collect(MethodFrame, "Method2", "", MethodAdvice)

	return "success", nil
}

func (s *sampleStruct) Method3(arg1 string, arg2 int) (string, error) {
	var err error
	ctx := Before(context.Background())
	defer func() {After(ctx, err)}()

	s.privateMethod1(arg1, arg2)

	return "success", nil
}

func (s *sampleStruct) Special(arg1 string, arg2 int) (string, error) {
	var err error
	ctx := Before(context.Background())
	defer func() {After(ctx, err)}()

	s.collector.Collect(MethodFrame, "Special", "", MethodAdvice)

	return "success", nil
}

func (s *sampleStruct) privateMethod1(arg1 string, arg2 int) (string, error) {
	var err error
	ctx := Before(context.Background())
	defer func() {After(ctx, err)}()

	s.collector.Collect(MethodFrame, "privateMethod1", "", MethodAdvice)

	return "success", nil
}


type loggingAspect struct {
	collector		*aspectCollector
}

func (l *loggingAspect) Before(ctx context.Context) context.Context {
	definition := AspectFromContext(ctx)
	l.collector.Collect(BeforeFrame, definition.MethodName, definition.CallingMethodName, LoggingAdvice)
	return ctx
}

func (l *loggingAspect) After(ctx context.Context, err error) context.Context {
	definition := AspectFromContext(ctx)
	l.collector.Collect(AfterFrame, definition.MethodName, definition.CallingMethodName, LoggingAdvice)
	return ctx
}

type countingAspect struct {
	collector		*aspectCollector
}

func (c *countingAspect) Before(ctx context.Context) context.Context {
	definition := AspectFromContext(ctx)
	c.collector.Collect(BeforeFrame, definition.MethodName, definition.CallingMethodName, CountAdvice)
	return ctx
}

func (c *countingAspect) After(ctx context.Context, err error) context.Context {
	definition := AspectFromContext(ctx)
	c.collector.Collect(AfterFrame, definition.MethodName, definition.CallingMethodName, CountAdvice)
	return ctx
}


type methodCall struct {
	frame					string
	methodName				string
	calingMethodName		string
	op						string
}

type aspectCollector struct {
	methodCalls		[]methodCall
}

func (a *aspectCollector) Collect(frame, methodName, callingMethodName, op string) *aspectCollector {
	a.methodCalls = append(a.methodCalls, methodCall{frame, methodName, callingMethodName, op})
	return a
}



