package middleware

import (
	"net/http"
	"warehouse-service/internal/util/liberr"
	"warehouse-service/internal/util/librest"
	"warehouse-service/module/warehouse/entity"
)

type basicAuthMiddleware struct {
	username string
	password string
	handler  librest.GatewayHandler
}

func WithBasicAuthMiddleware(username string, password string) librest.GatewayMiddleware {
	return func(handle librest.GatewayHandler) librest.GatewayHandler {
		em := basicAuthMiddleware{
			username: username,
			password: password,
			handler:  handle,
		}

		return em.handle
	}
}

func (ba *basicAuthMiddleware) handle(w http.ResponseWriter, r *http.Request) error {
	username, password, ok := r.BasicAuth()
	if !ok || username != ba.username || password != ba.password {
		return liberr.ResolveError(entity.ErrorForbidden)
	}

	return ba.handler(w, r)
}
