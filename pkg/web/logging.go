package web

import (
	"github.com/jfbramlett/go-aop/pkg/logging"
	"net/http"
)

const (
	endpoint = "endpoint"
	requestMethod = "requestMethod"
)

type LoggingMiddleware struct {
	method		string
}


func (l *LoggingMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		logger, reqCtx := logging.Named(l.method).
			WithField(endpoint, r.RequestURI).
			And(requestMethod, r.Method).
			Create(r.Context())

		logger.Info("request received")
		next.ServeHTTP(w, r.WithContext(reqCtx))

		logger.Info("request completed")
	})
}
