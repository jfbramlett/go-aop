package aop

import (
	"context"
	"github.com/golang-collections/collections/stack"
	"regexp"
	"runtime"
	"strings"
)

const (
	UnknownService 		= "Unknown"
	UnknownMethod 		= "Unknown"
	CallsBackToMethod 	= 2
)

type aopCtxKey struct{}

var myAopCtxKey = aopCtxKey{}


// Advice is the interface implemented to handle a cross-cutting concern
type Advice interface {
	Before(ctx context.Context) context.Context
	After(ctx context.Context, err error) context.Context
}

// AspectMgr is responsible for handling identifying and running our cross cutting concern
type AspectMgr interface {
	GetServiceName() string
	RegisterJoinPoint(pointcut Pointcut, advice Advice)
	Before(ctx context.Context, method string) context.Context
	After(ctx context.Context, err error) context.Context
}

var globalAspectMgr AspectMgr

// Aspect represents a specific invocation of a cross-cutting concern
type Aspect struct {
	MethodName        	string
	joinPoints         	[]joinPoint
}

// Pointcut defines how we determine if a given advice is relevant to the specified method
type Pointcut interface {
	Matches(method string) bool
}

type regexPointcut struct {
	pattern			string
}

func (r *regexPointcut) Matches(method string) bool {
	matches, err := regexp.MatchString(r.pattern, method)
	return err == nil && matches
}

// NewRegexPointcut returns a new Pointcut that uses regex pattern matching
func NewRegexPointcut(pattern string) Pointcut {
	return &regexPointcut{pattern: pattern}
}

type joinPoint struct {
	pointcut 		Pointcut
	advice        	Advice
}

type aspectMgr struct {
	joinPoints  []joinPoint
	serviceName string
}

// GetServiceName gets the name of the service we are running in
func (a *aspectMgr) GetServiceName() string {
	return a.serviceName
}

// RegisterJoinPoint registers a new advice for a given pointcut. The pointcut is a regex pattern used to match against a method name
func (a *aspectMgr) RegisterJoinPoint(pointcut Pointcut, advice Advice) {
	a.joinPoints = append(a.joinPoints, joinPoint{pointcut: pointcut, advice: advice})
}

// Before loops over all of the registered joinpoints and executes the Before advice for those whose pointcuts match
func (a *aspectMgr) Before(ctx context.Context, method string) context.Context {
	ac := &Aspect{joinPoints: make([]joinPoint, 0), MethodName: method}

	for _, k := range a.joinPoints {
		if k.pointcut.Matches(method) {
			ac.joinPoints = append(ac.joinPoints, k)
		}
	}

	if len(ac.joinPoints) > 0 {
		ctx = a.contextWithAspect(ctx, ac)

		for _, r := range ac.joinPoints {
			ctx = r.advice.Before(ctx)
		}
	}

	return ctx
}

func (a *aspectMgr) After(ctx context.Context, err error) context.Context {
	aop  := AspectFromContext(ctx)
	if aop != nil {
		for i := len(aop.joinPoints) - 1; i >= 0 ; i-- {
			ctx = aop.joinPoints[i].advice.After(ctx, err)
		}
	}

	return a.removeAspectFromContext(ctx)
}

func (a *aspectMgr) contextWithAspect(ctx context.Context, aop *Aspect) context.Context {
	return PushToContext(ctx, myAopCtxKey, aop)
}

func (a *aspectMgr) removeAspectFromContext(ctx context.Context) context.Context {
	ctx, _ = PopFromContext(ctx, myAopCtxKey)
	return ctx
}



func InitAOP(service string) {
	globalAspectMgr = &aspectMgr{serviceName: service, joinPoints: make([]joinPoint, 0)}
}

func GetServiceName() string {
	if globalAspectMgr != nil {
		return globalAspectMgr.GetServiceName()
	} else {
		return UnknownService
	}
}

// RegisterJoinPoint is function used to register a new advice with the given regex pointcut (will be compared
// against the calling method
func RegisterJoinPoint(pointcut Pointcut, advice Advice) {
	if globalAspectMgr != nil {
		globalAspectMgr.RegisterJoinPoint(pointcut, advice)
	}
}

// Before is the function invoked at the start of a method to execute any registered joinPoints
func Before(ctx context.Context) context.Context {
	if globalAspectMgr != nil {
		return globalAspectMgr.Before(ctx, getMethodNameAtOffset(CallsBackToMethod))
	} else {
		return ctx
	}
}

func After(ctx context.Context, err error) context.Context {
	if globalAspectMgr != nil {
		return globalAspectMgr.After(ctx, err)
	} else {
		return ctx
	}
}

func GetMethodName() string {
	return getMethodNameAtOffset(CallsBackToMethod)
}

func getMethodNameAtOffset(offset int) string {
	pc, _, _, ok := runtime.Caller(offset)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		return details.Name()
	}

	return UnknownMethod
}

func MethodNameFromFullPath(fullMethod string) string {
	idx := strings.LastIndex(fullMethod, ".")
	if idx > 0 {
		return fullMethod[idx+1:]
	}
	return fullMethod
}

func AspectFromContext(ctx context.Context) *Aspect {
	ctxStack := ctx.Value(myAopCtxKey)
	if ctxStack != nil {
		if aopStack, ok := ctxStack.(*stack.Stack); ok {
			stackItem := aopStack.Peek()
			if stackItem != nil {
				if aopItem, ok := stackItem.(*Aspect); ok {
					return aopItem
				}
			}
		}
	}

	return nil
}

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