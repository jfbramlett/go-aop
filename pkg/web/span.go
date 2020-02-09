package web

import (
	"github.com/jfbramlett/go-aop/pkg/tracing"
	"net/http"
)

type SpanMiddleware struct {
	method		string
}


func (s *SpanMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		sp, spanCtx := tracing.StartSpanFromContext(r.Context(), r.RequestURI)

		next.ServeHTTP(w, r.WithContext(spanCtx))

		sp.Finish()
	})
}
