package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"shop-service/internal/config"
	"syscall"
)

func main() {
	gatewayServer, cfg, err := config.NewGatewayServer()
	if err != nil {
		log.Fatalf("failed to create new gateway server: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func(server *http.Server) {
		log.Printf("starting gateway server on: %v\n", cfg.GatewayHost)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error listening gateway server: %v", err)
		}
	}(gatewayServer)

	<-sigChan

	log.Println("shutting down the gateway server...")
	err = gatewayServer.Shutdown(context.Background())
	if err != nil {
		log.Println(err)
	}
	log.Println("gateway server gracefully stopped")
}
