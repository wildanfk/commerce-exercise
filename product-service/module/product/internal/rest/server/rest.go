package server

import (
	"net/http"
	"product-service/internal/util/librest"
	"product-service/module/product/internal/rest/handler"
	"product-service/module/product/internal/rest/middleware"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type ServerConfig struct {
	Usecases *ServerUsecase
	Logger   *zap.Logger

	BasicAuthUsername string
	BasicAuthPassword string
}

type ServerUsecase struct {
	Product handler.ProductUsecase
}

func RegisterRESTHandler(serverMux *mux.Router, cfg *ServerConfig) error {
	shop := handler.NewProductHandler(cfg.Usecases.Product)

	registerInternalHandler(serverMux, cfg, http.MethodGet, "/check-products", shop.CheckProduct)

	registerHandler(serverMux, cfg, http.MethodGet, "/products", shop.ListProduct)

	return nil
}

func registerInternalHandler(serverMux *mux.Router, cfg *ServerConfig, method, path string, handle librest.GatewayHandler) {
	basicAuthMiddlewares := []librest.GatewayMiddleware{
		middleware.WithBasicAuthMiddleware(cfg.BasicAuthUsername, cfg.BasicAuthPassword),
		middleware.WithErrorMiddleware(),
		librest.WithLoggingMiddleware(path, cfg.Logger),
	}
	gatewayHandler := librest.GatewayHandlerFunc(librest.ApplyGatewayMiddlewares(handle, basicAuthMiddlewares...))

	serverMux.HandleFunc(path, gatewayHandler).Methods(method)
}

func registerHandler(serverMux *mux.Router, cfg *ServerConfig, method, path string, handle librest.GatewayHandler) {
	basicAuthMiddlewares := []librest.GatewayMiddleware{
		middleware.WithErrorMiddleware(),
		librest.WithLoggingMiddleware(path, cfg.Logger),
	}
	gatewayHandler := librest.GatewayHandlerFunc(librest.ApplyGatewayMiddlewares(handle, basicAuthMiddlewares...))

	serverMux.HandleFunc(path, gatewayHandler).Methods(method)
}
