package main

import (
	"github.com/kameikay/service-orchestration/configs"
	"github.com/kameikay/service-orchestration/internal/infra/web/controllers"
	"github.com/kameikay/service-orchestration/internal/infra/web/handlers"
	"github.com/kameikay/service-orchestration/internal/infra/web/webserver"
	"github.com/kameikay/service-orchestration/internal/service"
)

func main() {
	err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	server := webserver.NewWebServer(":8081")

	server.MountMiddlewares()

	viaCepService := service.NewViaCepService()
	weatherApiService := service.NewWeatherApiService()
	handler := handlers.NewHandler(viaCepService, weatherApiService)
	controller := controllers.NewController(server.Router, handler)
	controller.Route()

	server.Start()
}
