package librest

import (
	"net/http"
)

// GatewayHandler define handler signature for gateway middleware
type GatewayHandler func(http.ResponseWriter, *http.Request) error

// GatewayMiddleware define middleware of GatewayHandler
type GatewayMiddleware func(GatewayHandler) GatewayHandler

// ApplyGatewayMiddlewares return GatewayHandler with Middlewares.
func ApplyGatewayMiddlewares(handle GatewayHandler, ds ...GatewayMiddleware) GatewayHandler {
	for _, d := range ds {
		handle = d(handle)
	}
	return handle
}

func GatewayHandlerFunc(h GatewayHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if x, ok := w.(*ResponseWriter); ok {
			x.SetError(err)
		}
	}
}
