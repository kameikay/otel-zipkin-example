package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/kameikay/service-input/configs"
	"github.com/kameikay/service-input/internal/infra/web/controllers"
	"github.com/kameikay/service-input/internal/infra/web/handlers"
	"github.com/kameikay/service-input/internal/infra/web/webserver"
	"github.com/kameikay/service-input/internal/service"
)

func main() {
	err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	signChannel := make(chan os.Signal, 1)
	signal.Notify(signChannel, os.Interrupt)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	shutdown, err := configs.SetupOTel(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown tracer provider: ", err)
		}
	}()

	server := webserver.NewWebServer(":8080")
	server.MountMiddlewares()

	apiService := service.NewGetTemperatureService()
	handler := handlers.NewHandler(apiService)
	controller := controllers.NewController(server.Router, handler)
	controller.Route()

	go func() {
		server.Start()
	}()

	select {
	case <-signChannel:
		log.Println("shutting down server gracefully...")
	case <-ctx.Done():
		log.Println("shutting down server...")
	}

	_, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
}
