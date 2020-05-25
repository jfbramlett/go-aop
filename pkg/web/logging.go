package web

import (
	"net/http"

	"github.com/jfbramlett/go-aop/pkg/tracing"

	"github.com/jfbramlett/go-aop/pkg/logging"
)

const (
	endpoint      = "endpoint"
	requestMethod = "requestMethod"
	traceId       = "traceId"
)

type LoggingMiddleware struct {
	method string
}

func (l *LoggingMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		requestId := tracing.GetTraceFromContext(r.Context())

		logger, reqCtx := logging.Named(l.method).
			WithField(endpoint, r.RequestURI).
			And(requestMethod, r.Method).
			And(traceId, requestId).
			Create(r.Context())

		logger.Info("request received")
		next.ServeHTTP(w, r.WithContext(reqCtx))

		logger.Info("request completed")
	})
}
