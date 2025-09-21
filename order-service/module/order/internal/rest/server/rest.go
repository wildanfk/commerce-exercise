package server

import (
	"net/http"
	"order-service/internal/util/librest"
	"order-service/module/order/internal/rest/handler"
	"order-service/module/order/internal/rest/middleware"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type ServerConfig struct {
	Usecases *ServerUsecase
	Logger   *zap.Logger

	AuthServiceJWTSecret string
}

type ServerUsecase struct {
	Order handler.OrderUsecase
}

func RegisterRESTHandler(serverMux *mux.Router, cfg *ServerConfig) error {
	order := handler.NewOrderHandler(cfg.Usecases.Order, handler.OrderHandlerConfig{
		AuthServiceJWTSecret: cfg.AuthServiceJWTSecret,
	})

	registerHandler(serverMux, cfg, http.MethodPost, "/checkout-orders", order.CreateOrder)

	return nil
}

func registerHandler(serverMux *mux.Router, cfg *ServerConfig, method, path string, handle librest.GatewayHandler) {
	basicAuthMiddlewares := []librest.GatewayMiddleware{
		middleware.WithErrorMiddleware(),
		librest.WithLoggingMiddleware(path, cfg.Logger),
	}
	gatewayHandler := librest.GatewayHandlerFunc(librest.ApplyGatewayMiddlewares(handle, basicAuthMiddlewares...))

	serverMux.HandleFunc(path, gatewayHandler).Methods(method)
}
