package messaging

// MiddlewareFunc is a function which receives an http.Handler and returns another http.Handler.
// Typically, the returned handler is a closure which does something with the http.ResponseWriter and http.Request passed
// to it, and then calls the handler passed as parameter to the MiddlewareFunc.
type MiddlewareFunc func(handlerFunc OutboxHandlerFunc) OutboxHandlerFunc

// middleware interface is anything which implements a MiddlewareFunc named Middleware.
type middleware interface {
	Middleware(next OutboxHandlerFunc) OutboxHandlerFunc
}

// Middleware allows MiddlewareFunc to implement the middleware interface.
func (mw MiddlewareFunc) Middleware(next OutboxHandlerFunc) OutboxHandlerFunc {
	return mw(next)
}
