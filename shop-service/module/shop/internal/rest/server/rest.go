package server

import (
	"net/http"
	"shop-service/internal/util/librest"
	"shop-service/module/shop/internal/rest/handler"
	"shop-service/module/shop/internal/rest/middleware"

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
	Shop handler.ShopUsecase
}

func RegisterRESTHandler(serverMux *mux.Router, cfg *ServerConfig) error {
	shop := handler.NewShopHandler(cfg.Usecases.Shop)

	registerInternalHandler(serverMux, cfg, http.MethodGet, "/shops", shop.ListByParams)

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
