package config

import (
	"order-service/module/order/internal/rest/server"

	"github.com/gorilla/mux"
)

func RegisterGatewayHandler(serverMux *mux.Router, cfg *OrderConfig) error {
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
			Order: usecases.orderUsecase,
		},
		Logger:               cfg.Logger,
		AuthServiceJWTSecret: cfg.AuthServiceJWTSecret,
	}
	if err := server.RegisterRESTHandler(serverMux, serverConfig); err != nil {
		return err
	}

	return nil
}
