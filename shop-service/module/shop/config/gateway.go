package config

import (
	"shop-service/module/shop/internal/rest/server"

	"github.com/gorilla/mux"
)

func RegisterGatewayHandler(serverMux *mux.Router, cfg *ShopConfig) error {
	repositories, err := newRepositories(cfg)
	if err != nil {
		return err
	}

	usecases, err := newUsecase(repositories)
	if err != nil {
		return err
	}

	serverConfig := &server.ServerConfig{
		Usecases: &server.ServerUsecase{
			Shop: usecases.shopUsecase,
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
