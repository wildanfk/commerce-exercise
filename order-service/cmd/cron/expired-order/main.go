package main

import (
	"log"
	"order-service/internal/config"
)

func main() {
	cron, err := config.NewCronExpiredOrder()
	if err != nil {
		log.Fatalf("failed to create new cron job: %v", err)
	}

	err = cron.ExecuteCron()
	if err != nil {
		log.Printf("execute cron job error = %v\n", err)
	}

	log.Println("shutting down the cron job")
	log.Println("cron job gracefully stopped")
}
