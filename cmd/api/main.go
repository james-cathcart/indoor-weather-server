package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"weatherserver/internal/system"
	"weatherserver/internal/weather"
)

var (
	ElasticHost       string
	WeatherServerPort string
)

func init() {
	ElasticHost = os.Getenv(`ELASTIC_HOST`)
	WeatherServerPort = os.Getenv(`WEATHER_PORT`)
	if WeatherServerPort == `` {
		WeatherServerPort = `8080`
	}
}

func main() {

	weatherService := weather.NewElasticService(&http.Client{}, ElasticHost)
	weatherHandler := weather.NewAPI(weatherService)

	mux := http.NewServeMux()

	mux.Handle(`/ingest/weather`, weatherHandler)
	mux.HandleFunc(`/system/health`, system.HealthCheck)

	server := http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", WeatherServerPort),
		Handler: mux,
	}

	log.Println(`starting server...`)

	err := server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}

	log.Println(`application exiting`)
}
