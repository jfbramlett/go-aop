package web

import (
	"net/http"

	"github.com/jfbramlett/go-aop/pkg/tracing"

	"github.com/google/uuid"
)

const (
	headerRequestId = "xxx-request-id"
)

type TracingMiddleware struct {
	method string
}

func (l *TracingMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		requestId := r.Header.Get(headerRequestId)
		if requestId == "" {
			requestId = uuid.New().String()
		}

		reqCtx := tracing.SetTraceInContext(r.Context(), requestId)

		next.ServeHTTP(w, r.WithContext(reqCtx))
	})
}
