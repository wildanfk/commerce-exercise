package server

import (
	"net/http"
	"user-service/internal/util/librest"
	"user-service/module/auth/internal/rest/handler"
	"user-service/module/auth/internal/rest/middleware"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type ServerConfig struct {
	Usecases *ServerUsecase
	Logger   *zap.Logger
}

type ServerUsecase struct {
	Auth handler.AuthUsecase
}

func RegisterRESTHandler(serverMux *mux.Router, cfg *ServerConfig) error {
	auth := handler.NewAuthHandler(cfg.Usecases.Auth)

	registerHandler(serverMux, cfg, http.MethodPost, "/authentication", auth.Authentication)

	return nil
}

func registerHandler(serverMux *mux.Router, cfg *ServerConfig, method, path string, handle librest.GatewayHandler) {
	middlewares := []librest.GatewayMiddleware{
		middleware.WithErrorMiddleware(),
		librest.WithLoggingMiddleware(path, cfg.Logger),
	}
	gatewayHandler := librest.GatewayHandlerFunc(librest.ApplyGatewayMiddlewares(handle, middlewares...))

	serverMux.HandleFunc(path, gatewayHandler).Methods(method)
}
