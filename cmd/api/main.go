package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"weatherserver/internal/weather"
)

var (
	ElasticHost string
)

func init() {
	ElasticHost = os.Getenv(`ELASTIC_HOST`)
}
func main() {

	weatherService := weather.NewService(&http.Client{}, ElasticHost)
	weatherHandler := weather.NewAPI(weatherService)

	mux := http.NewServeMux()

	mux.Handle(`/ingest/weather`, weatherHandler)

	server := http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", 8080),
		Handler: mux,
	}

	log.Println(`starting server...`)

	err := server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}

	log.Println(`application exiting`)
}
