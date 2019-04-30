package aop

import (
	"context"
	"fmt"
	"testing"
)

func TestAOP(t *testing.T) {
	InitAOP("testAop")
	RegisterAspect(".*", &LoggingAspect{})
	RegisterAspect(".*Method1$", &CountingAspect{})

	st := SampleStruct{}
	st.Method1("arg1", 1)
	st.Method2("arg1", 1)
	st.Method3("arg1", 1)
}


type SampleStruct struct {

}

func (s *SampleStruct) Method1(arg1 string, arg2 int) (string, error) {
	var err error
	ctx := Before(context.Background())
	defer func() {After(ctx, err)}()

	fmt.Println("In method 1")
	return "success", nil
}

func (s *SampleStruct) Method2(arg1 string, arg2 int) (string, error) {
	var err error
	ctx := Before(context.Background())
	defer func() {After(ctx, err)}()

	fmt.Println("In method 2")
	return "success", nil
}

func (s *SampleStruct) Method3(arg1 string, arg2 int) (string, error) {
	var err error
	ctx := Before(context.Background())
	defer func() {After(ctx, err)}()

	s.privateMethod1(arg1, arg2)

	fmt.Println("In method 3")
	return "success", nil
}

func (s *SampleStruct) privateMethod1(arg1 string, arg2 int) (string, error) {
	var err error
	ctx := Before(context.Background())
	defer func() {After(ctx, err)}()

	fmt.Println("In private method 1")
	return "success", nil
}


type LoggingAspect struct {

}

func (l *LoggingAspect) Before(ctx context.Context) context.Context {
	definition := AOPFromContext(ctx)
	fmt.Println("Executing logging method before " + definition.MethodName + " and " + definition.CallingMethodName)
	return ctx
}

func (l *LoggingAspect) After(ctx context.Context, err error) context.Context {
	definition := AOPFromContext(ctx)
	fmt.Println("Executing logging method after " + definition.MethodName + " and " + definition.CallingMethodName)
	return ctx
}

type CountingAspect struct {

}

func (l *CountingAspect) Before(ctx context.Context) context.Context {
	definition := AOPFromContext(ctx)
	fmt.Println("Executing count method before " + definition.MethodName + " and " + definition.CallingMethodName)
	return ctx
}

func (l *CountingAspect) After(ctx context.Context, err error) context.Context {
	definition := AOPFromContext(ctx)
	fmt.Println("Executing count method after " + definition.MethodName + " and " + definition.CallingMethodName)
	return ctx
}
