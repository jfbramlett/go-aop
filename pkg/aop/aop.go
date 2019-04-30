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


// Advice is the implementation of the logic to perform around a given method
type Advice interface {
	Before(ctx context.Context) context.Context
	After(ctx context.Context, err error) context.Context
}

type joinPoint struct {
	pointcut 		string
	advice        	Advice
}

type aspectMgr struct {
	joinPoints  []joinPoint
	serviceName string
}

type Aspect struct {
	MethodName        	string
	CallingMethodName 	string
	joinPoints         	[]joinPoint
}

var globalAspectMgr = aspectMgr{joinPoints: make([]joinPoint, 0), serviceName: "unknown"}

func InitAOP(service string) {
	globalAspectMgr = aspectMgr{joinPoints: make([]joinPoint, 0), serviceName: service}
}

func GetServiceName() string {
	return globalAspectMgr.serviceName
}

// RegisterJoinPoint is function used to register a new advice with the given regex pointcut (will be compared
// against the calling method
func RegisterJoinPoint(pointcut string, advice Advice) {
	globalAspectMgr.joinPoints = append(globalAspectMgr.joinPoints, joinPoint{pointcut: pointcut, advice: advice})
}

// Before is the function invoked at the start of a method to execute any registered joinPoints
func Before(ctx context.Context) context.Context {
	fullMethod := getMethod()
	ac := &Aspect{joinPoints: make([]joinPoint, 0), CallingMethodName: getCallingMethod(),
		MethodName: methodNameFromFullPath(fullMethod)}

	for _, k := range globalAspectMgr.joinPoints {
		matches, err := regexp.MatchString(k.pointcut, fullMethod)
		if err == nil && matches {
			ac.joinPoints = append(ac.joinPoints, k)
		}
	}

	if len(ac.joinPoints) > 0 {
		ctx = contextWithAspect(ctx, ac)

		for _, r := range ac.joinPoints {
			ctx = r.advice.Before(ctx)
		}
	}

	return ctx
}

func After(ctx context.Context, err error) context.Context {
	aop  := AspectFromContext(ctx)
	if aop != nil {
		for i := len(aop.joinPoints) - 1; i >= 0 ; i-- {
			ctx = aop.joinPoints[i].advice.After(ctx, err)
		}
	}

	return removeAspectFromContext(ctx)
}

func contextWithAspect(ctx context.Context, aop *Aspect) context.Context {
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

func removeAspectFromContext(ctx context.Context) context.Context {
	ctxStack := ctx.Value(myAopCtxKey)
	if ctxStack != nil {
		aopStack, valid := ctxStack.(*stack.Stack)
		if valid {
			aopStack.Pop()
		}
	}

	return ctx
}

func AspectFromContext(ctx context.Context) *Aspect {
	ctxStack := ctx.Value(myAopCtxKey)
	if ctxStack != nil {
		aopStack, valid := ctxStack.(*stack.Stack)
		if valid {
			stackItem := aopStack.Peek()
			if stackItem != nil {
				aopItem, valid := stackItem.(*Aspect)
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
		return details.Name()
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

