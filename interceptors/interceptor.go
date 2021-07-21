package interceptors

import (
	"net/http"
)

type MiddlewareInterceptor func(http.ResponseWriter, *http.Request, http.HandlerFunc)

type MiddlewareHandlerFunc http.HandlerFunc

type MiddlewareHandler http.Handler

func (cont MiddlewareHandlerFunc) Intercept(mw MiddlewareInterceptor) MiddlewareHandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		mw(writer, request, http.HandlerFunc(cont))
	}
}

type MiddlewareChain []MiddlewareInterceptor

func (chain MiddlewareChain) Handler(handler http.HandlerFunc) MiddlewareHandler {
	curr := MiddlewareHandlerFunc(handler)
	for i := len(chain) - 1; i >= 0; i-- {
		mw := chain[i]
		curr = curr.Intercept(mw)
	}
	return http.HandlerFunc(curr)
}
