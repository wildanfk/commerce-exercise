package config

import (
	"warehouse-service/module/warehouse/internal/rest/server"

	"github.com/gorilla/mux"
)

func RegisterGatewayHandler(serverMux *mux.Router, cfg *WarehouseConfig) error {
	repositories, err := newRepositories(cfg)
	if err != nil {
		return err
	}

	usecases, err := newUsecase(cfg, repositories)
	if err != nil {
		return err
	}

	serverConfig := &server.ServerConfig{
		Usecases: &server.ServerUsecase{
			Warehouse:      usecases.warehouseUsecase,
			WarehouseStock: usecases.warehouseStockUsecase,
		},
		Logger:            cfg.Logger,
		BasicAuthUsername: cfg.BasicAuthUsername,
		BasicAuthPassword: cfg.BasicAuthPassword,
	}
	if err := server.RegisterRESTHandler(serverMux, serverConfig); err != nil {
		return err
	}

	return nil
}
