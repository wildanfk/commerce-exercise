package config

import (
	"user-service/module/auth/internal/rest/server"

	"github.com/gorilla/mux"
)

func RegisterGatewayHandler(serverMux *mux.Router, cfg *AuthConfig) error {
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
			Auth: usecases.authUsecase,
		},
		Logger: cfg.Logger,
	}
	if err := server.RegisterRESTHandler(serverMux, serverConfig); err != nil {
		return err
	}

	return nil
}
