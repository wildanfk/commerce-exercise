package server

import (
	"net/http"
	"warehouse-service/internal/util/librest"
	"warehouse-service/module/warehouse/internal/rest/handler"
	"warehouse-service/module/warehouse/internal/rest/middleware"

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
	Warehouse      handler.WarehouseUsecase
	WarehouseStock handler.WarehouseStockUsecase
}

func RegisterRESTHandler(serverMux *mux.Router, cfg *ServerConfig) error {
	warehouse := handler.NewWarehouseHandler(cfg.Usecases.Warehouse)

	registerInternalHandler(serverMux, cfg, http.MethodPost, "/warehouse-actives", warehouse.WarehouseActivation)

	warehouseStock := handler.NewWarehouseStockHandler(cfg.Usecases.WarehouseStock)

	registerInternalHandler(serverMux, cfg, http.MethodGet, "/active-stocks", warehouseStock.ActiveStock)
	registerInternalHandler(serverMux, cfg, http.MethodPost, "/adjustment-stocks", warehouseStock.AdjustmentStock)
	registerInternalHandler(serverMux, cfg, http.MethodPost, "/transfer-stocks", warehouseStock.TransferStock)

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
