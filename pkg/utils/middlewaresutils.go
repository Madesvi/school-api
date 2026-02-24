package utils

import "net/http"

// Middleware - like a check point который стоит между клиентским запросом
// и Запросом к серверу он может инспектировать, изменять или блокировать
// request перед приемом
// И также делает тоже самое с responsons к Client
// Data validation

// Middleware is a function that wraps an http.Hadler with addition functionality
type Middleware func(http.Handler) http.Handler

func ApplyMiddleWares(handler http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}
