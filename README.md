# Open Telemetry and Zipkin - Go

## Objective
The goal of this project is to develop a system in Go that takes a postal code as input, identifies the city, and returns the current weather (temperature in Celsius, Fahrenheit, and Kelvin), using Open Telemetry and Zipkin to trace the requests.

## Technologies and Tools Used

- Programming Language: Go
- External APIs: viaCEP and WeatherAPI
- Observability: Open Telemetry and Zipkin
- Containerization: Docker


## Prerequisites

- Go 1.16+ installed
- Docker installed (for running the project via Docker)
- An internet connection to access external APIs (viaCEP and WeatherAPI)


## How to run the project

### Environment Variables

- Service Input:

-- WEATHER_SERVICE_URL = http://service-orchestration:8081/
-- SERVICE_NAME = service-input
-- OTEL_COLLECTOR_ADDR = otel-collector:4317

- Service Orchestration:

-- WEATHER_API_KEY = {YOUR_API_KEY} (WeatherAPI)[https://www.weatherapi.com/]
-- SERVICE_NAME = service-orchestration
-- OTEL_COLLECTOR_ADDR = otel-collector:4317

### Running via docker-file

1. Run the following command to start the application:
```bash
docker-compose up
```
2. Make a request to the service-input:
```bash
curl --request POST --url 'http://localhost:8080' -H "Content-Type: application/json" -d '{"cep" : "01001-000"}'
```

## Zipkin

To access the Zipkin dashboard, open your browser and go to the following address:
```bash
http://localhost:9411
```