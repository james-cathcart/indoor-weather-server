package main

import (
	"fmt"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	stdlog "log"
	"net/http"
	"os"
	"weatherserver/internal/common"
	"weatherserver/internal/logging"
	"weatherserver/internal/system"
	"weatherserver/internal/weather"
)

const (
	logFileName string = `log.json`
)

var (
	ElasticHost       string
	WeatherServerPort string
	Environment       string
	LogLevel          string
)

func init() {

	ElasticHost = os.Getenv(common.ElasticHost)
	WeatherServerPort = os.Getenv(common.WeatherPort)
	if WeatherServerPort == `` {
		WeatherServerPort = `8080`
	}

	// determine deployment environment
	env, ok := os.LookupEnv(common.WeatherEnv)
	if !ok {
		stdlog.Fatalf("error: could not determine deployment environment")
	}
	Environment = env
	LogLevel, ok = os.LookupEnv(common.LogLevel)
	if !ok {
		stdlog.Println("warning: could not determine log level from environment, applying default `error` level")
		LogLevel = common.LogError
	}
}

func main() {

	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		stdlog.Fatal(err)
	}
	defer func(closeFunc func() error) {
		err = logFile.Close()
		if err != nil {
			stdlog.Println(err)
		}
	}(logFile.Close)

	encoderConfig := ecszap.NewDefaultEncoderConfig()

	var zapCoreLevel zapcore.Level
	switch LogLevel {
	case `info`:
		zapCoreLevel = zap.InfoLevel
	case `debug`:
		zapCoreLevel = zap.DebugLevel
	case `warn`:
		zapCoreLevel = zap.WarnLevel
	case `error`:
		zapCoreLevel = zap.ErrorLevel
	}

	logCore := ecszap.NewCore(encoderConfig, logFile, zapCoreLevel)
	zapLog := zap.New(logCore, zap.AddCaller())
	zapLog = zapLog.Named(`weather-server`)
	zapLog = zapLog.With(zap.String(`env`, Environment))

	customLogger := logging.NewElasticLogger(zapLog)

	weatherService := weather.NewElasticService(&http.Client{}, ElasticHost, customLogger)
	weatherHandler := weather.NewAPI(weatherService, customLogger)

	systemHandler := system.NewSystemHandler(customLogger)

	mux := http.NewServeMux()

	mux.Handle(`/ingest/weather`, weatherHandler)
	mux.Handle(`/system`, systemHandler)

	server := http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", WeatherServerPort),
		Handler: mux,
	}

	stdlog.Println(`starting server...`)

	err = server.ListenAndServe()
	if err != nil {
		stdlog.Printf("%s: error: %v", Environment, err)
	}

	stdlog.Println(`application exiting`)
}
