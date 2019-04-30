package aop

import (
	"context"
	"github.com/golang-collections/collections/stack"
	"regexp"
	"runtime"
	"strings"
)

const UnknownMethod = "Unknown"

type aopCtxKey struct{}

var myAopCtxKey = aopCtxKey{}



// Aspect is the implementation of the logic to perform around a given method
type Aspect interface {
	Before(ctx context.Context) context.Context
	After(ctx context.Context, err error) context.Context
}

type aopWrapper struct {
	methodPattern	string
	aspect 			Aspect
}

type Aop struct {
	Service           	string
	MethodName        	string
	CallingMethodName 	string
	aspects           	[]aopWrapper
}

var globalAspect = Aop{aspects: make([]aopWrapper, 0), Service: "unknown"}

func InitAOP(service string) {
	globalAspect.Service = service
}

// RegisterAspect is function used to register a new aspect, it wraps a method which is matched via regex
func RegisterAspect(methodPattern string, aspect Aspect) {
	globalAspect.aspects = append(globalAspect.aspects, aopWrapper{methodPattern: methodPattern, aspect: aspect})
}

// Before is the function invoked at the start of a method to execute any registered aspects
func Before(ctx context.Context) context.Context {
	method := getMethod()
	ac := &Aop{Service: globalAspect.Service, aspects: make([]aopWrapper, 0), CallingMethodName: getCallingMethod(),
		MethodName: method}

	for _, k := range globalAspect.aspects {
		matches, err := regexp.MatchString(k.methodPattern, method)
		if err == nil && matches {
			ac.aspects = append(ac.aspects, k)
		}
	}

	if len(ac.aspects) > 0 {
		ctx = contextWithAOP(ctx, ac)

		for _, r := range ac.aspects {
			ctx = r.aspect.Before(ctx)
		}
	}

	return ctx
}

func After(ctx context.Context, err error) context.Context {
	aop  := AOPFromContext(ctx)
	if aop != nil {
		for i := len(aop.aspects) - 1; i >= 0 ; i-- {
			ctx = aop.aspects[i].aspect.After(ctx, err)
		}
	}

	return removeAOPFromContext(ctx)
}

func contextWithAOP(ctx context.Context, aop *Aop) context.Context {
	var aopStack *stack.Stack
	var valid bool

	ctxStack := ctx.Value(myAopCtxKey)
	if ctxStack == nil {
		aopStack = stack.New()
		ctx = context.WithValue(ctx, myAopCtxKey, aopStack)
	} else {
		aopStack, valid = ctxStack.(*stack.Stack)
		if !valid {
			return ctx
		}
	}

	aopStack.Push(aop)

	return ctx
}

func removeAOPFromContext(ctx context.Context) context.Context {
	ctxStack := ctx.Value(myAopCtxKey)
	if ctxStack != nil {
		aopStack, valid := ctxStack.(*stack.Stack)
		if valid {
			aopStack.Pop()
		}
	}

	return ctx
}

func AOPFromContext(ctx context.Context) *Aop {
	ctxStack := ctx.Value(myAopCtxKey)
	if ctxStack != nil {
		aopStack, valid := ctxStack.(*stack.Stack)
		if valid {
			stackItem := aopStack.Peek()
			if stackItem != nil {
				aopItem, valid := stackItem.(*Aop)
				if valid {
					return aopItem
				}
			}
		}
	}

	return nil
}

func getCallingMethod() string {
	pc, _, _, ok := runtime.Caller(3)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		return methodNameFromFullPath(details.Name())
	}

	return UnknownMethod
}

func getMethod() string {
	pc, _, _, ok := runtime.Caller(2)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		return methodNameFromFullPath(details.Name())
	}

	return UnknownMethod
}

func methodNameFromFullPath(fullMethod string) string {
	idx := strings.LastIndex(fullMethod, ".")
	if idx > 0 {
		return fullMethod[idx+1:]
	}
	return fullMethod
}

