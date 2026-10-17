package main

import (
	"log"

	"github.com/tousart/factory_analyst/app"
	"github.com/tousart/factory_analyst/config"
	"github.com/tousart/factory_analyst/pkg/server"
)

func main() {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting app in %s environment", cfg.Environment)

	if err := server.NewServer(cfg.Host(), cfg.GatewayHost(), app.NewApp().Transports...).Run(); err != nil {
		log.Fatal("Error starting app: ", err)
	}
}
