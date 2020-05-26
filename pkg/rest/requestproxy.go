package rest

import (
	"context"
	"net/http"

	"github.com/jfbramlett/go-aop/pkg/logging"

	"github.com/jfbramlett/go-aop/pkg/tracing"
	"github.com/jfbramlett/go-aop/pkg/web"
)

type RequestProxy interface {
	Before(ctx context.Context, r *http.Request) (*http.Request, error)
	After(ctx context.Context, r *http.Response, err error) (*http.Response, error)
}

type BaseRequestProxy struct {
}

func (b *BaseRequestProxy) Before(_ context.Context, r *http.Request) (*http.Request, error) {
	return r, nil
}

func (b *BaseRequestProxy) After(_ context.Context, r *http.Response, err error) (*http.Response, error) {
	return r, err
}

func NewTraceRequestProxy() RequestProxy {
	return &TraceRequestProxy{BaseRequestProxy: BaseRequestProxy{}}
}

type TraceRequestProxy struct {
	BaseRequestProxy
}

func (t *TraceRequestProxy) Before(ctx context.Context, r *http.Request) (*http.Request, error) {
	requestId := tracing.GetTraceFromContext(ctx)
	r.Header.Add(web.HeaderRequestId, requestId)
	return r, nil
}

func NewLoggingRequestProxy() RequestProxy {
	return &LoggingRequestProxy{BaseRequestProxy: BaseRequestProxy{}}
}

type LoggingRequestProxy struct {
	BaseRequestProxy
}

func (l *LoggingRequestProxy) Before(ctx context.Context, r *http.Request) (*http.Request, error) {
	logger := logging.GetLogger(ctx)
	logger.Infof("start executing %v for %v", r.Method, r.RequestURI)
	return r, nil
}

func (l *LoggingRequestProxy) After(ctx context.Context, r *http.Response, err error) (*http.Response, error) {
	logger := logging.GetLogger(ctx)
	if err != nil {
		logger.Error(err, "rest call failed")
	} else {
		logger.Info("call completed successfully")
	}
	return r, err
}
