package main

import (
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

	server := webserver.NewWebServer(":8080")

	server.MountMiddlewares()

	apiService := service.NewGetTemperatureService()
	handler := handlers.NewHandler(apiService)
	controller := controllers.NewController(server.Router, handler)
	controller.Route()

	server.Start()
}
