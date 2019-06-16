package aop

import (
	"context"
	"regexp"
	"runtime"
)

const (
	// UnknownService is a flag indicating the service is not known
	UnknownService    = "Unknown"
	// UnknownMethod represents an unknown method
	UnknownMethod     = "Unknown"
	// CallsBackToMethod indicates the number of steps back the call stack to use
	CallsBackToMethod = 2
)

type contextKey struct{}

var aopCtxKey = contextKey{}


// Advice is the interface implemented to handle a cross-cutting concern
type Advice interface {
	Before(ctx context.Context) context.Context
	After(ctx context.Context, err error)
}

// AspectMgr is responsible for handling identifying and running our cross cutting concern
type AspectMgr interface {
	GetServiceName() string
	RegisterJoinPoint(pointcut Pointcut, advice Advice)
	Before(ctx context.Context, method string) context.Context
	After(ctx context.Context, err error)
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
		ctx = context.WithValue(ctx, aopCtxKey, ac)

		for _, r := range ac.joinPoints {
			ctx = r.advice.Before(ctx)
		}
	}

	return ctx
}

func (a *aspectMgr) After(ctx context.Context, err error) {
	aop  := AspectFromContext(ctx)
	if aop != nil {
		for i := len(aop.joinPoints) - 1; i >= 0 ; i-- {
			aop.joinPoints[i].advice.After(ctx, err)
		}
	}
}

// InitAOP initializes our aspects
func InitAOP(service string) {
	globalAspectMgr = &aspectMgr{serviceName: service, joinPoints: make([]joinPoint, 0)}
}

// GetServiceName gets the name of the service
func GetServiceName() string {
	if globalAspectMgr != nil {
		return globalAspectMgr.GetServiceName()
	}
	return UnknownService
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
	}
	return ctx
}

// After is a global func used to execute our aspect
func After(ctx context.Context, err error) {
	if globalAspectMgr != nil {
		globalAspectMgr.After(ctx, err)
	}
}

// GetMethodName returns the method that has invoked the code that has invokved this method
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

// AspectFromContext gets the current aspect from the context
func AspectFromContext(ctx context.Context) *Aspect {
	ctxVal := ctx.Value(aopCtxKey)
	if ctxVal != nil {
		if aopItem, ok := ctxVal.(*Aspect); ok {
			return aopItem
		}
	}

	return nil
}