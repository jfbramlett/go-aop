package common

import (
	"context"
	"github.com/golang-collections/collections/stack"
	"runtime"
)

func PushToContext(ctx context.Context, ctxKey interface{}, value interface{}) context.Context {
	var dataStack *stack.Stack
	var ok bool

	ctxStack := ctx.Value(ctxKey)
	if ctxStack == nil {
		dataStack = stack.New()
		ctx = context.WithValue(ctx, ctxKey, dataStack)
	} else {
		if dataStack, ok = ctxStack.(*stack.Stack); !ok {
			return ctx
		}
	}

	dataStack.Push(value)

	return ctx
}

func PopFromContext(ctx context.Context, ctxKey interface{}) (context.Context, interface{}) {
	var item interface{}

	ctxStack := ctx.Value(ctxKey)
	if ctxStack != nil {
		if dataStack, ok := ctxStack.(*stack.Stack); ok {
			item = dataStack.Pop()
		}
	}

	return ctx, item
}

func FromContext(ctx context.Context, ctxKey interface{}) interface{} {
	ctxStack := ctx.Value(ctxKey)
	if ctxStack != nil {
		if aopStack, ok := ctxStack.(*stack.Stack); ok {
			return aopStack.Peek()
		}
	}

	return nil
}

func GetCallingMethodName() string {
	return GetMethodNameAt(3)
}

func GetMethodNameAt(idx int) string {
	pc, _, _, ok := runtime.Caller(idx)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		return details.Name()
	}

	return "unknown"
}
