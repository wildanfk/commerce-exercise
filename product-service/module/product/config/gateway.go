package config

import (
	"product-service/module/product/internal/rest/server"

	"github.com/gorilla/mux"
)

func RegisterGatewayHandler(serverMux *mux.Router, cfg *ProductConfig) error {
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
			Product: usecases.productUsecase,
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
