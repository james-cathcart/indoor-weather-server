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
	if ElasticHost == `` {
		stdlog.Fatalf("error: elastic host value empty")
	}

	WeatherServerPort = `8080`

	// determine deployment environment
	Environment = os.Getenv(common.WeatherEnv)
	if Environment == `` {
		stdlog.Fatal("error: could not determine deployment environment")
	}

	LogLevel = os.Getenv(common.LogLevel)
	if LogLevel == `` {
		stdlog.Println("warning: could not determine log level from environment, applying default `error` level")
		LogLevel = common.LogError
	}

	weather.ElasticIndex = os.Getenv(common.ElasticIndex)
	if weather.ElasticIndex == `` {
		stdlog.Fatal(`no elastic index found!`)
	}
}

func main() {

	/*****************************************************************

	The following section sets up logging for Elastic's Filebeat
	log aggregation technology.

	*****************************************************************/

	// create file for log data
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

	// create encoder to format logs into elastic search friendly format
	encoderConfig := ecszap.NewDefaultEncoderConfig()

	// map the log level
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

	// create the logger core
	logCore := ecszap.NewCore(encoderConfig, logFile, zapCoreLevel)

	// create and configure the new logger
	zapLog := zap.New(logCore, zap.AddCaller())
	zapLog = zapLog.Named(`weather-server`)
	zapLog = zapLog.With(zap.String(`env`, Environment))

	// create a new customer logger with the Elastic implementation.
	// the custom logger abstraction is logger agnostic and an
	// implementation can be written for any logger.
	customLogger := logging.NewElasticLogger(zapLog)

	/******************************************************************

	The following section handles the Dependency Injection for the
	application.

	******************************************************************/

	// create the weather service with the Elastic implementation (alternate
	// implementations can be written for your preferred database solution)
	weatherService := weather.NewElasticService(&http.Client{}, ElasticHost, customLogger)

	// create the weather REST API handler
	weatherHandler := weather.NewAPI(weatherService, customLogger)

	// create the system REST API handler
	systemHandler := system.NewSystemHandler(customLogger)

	/******************************************************************

	The following section creates and configures the HTTP server.

	******************************************************************/

	// create a new multiplexer
	mux := http.NewServeMux()

	// assign URIs to handlers
	mux.Handle(`/ingest/weather`, weatherHandler)
	mux.Handle(`/system`, systemHandler)

	// create/configure the server struct
	server := http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", WeatherServerPort),
		Handler: mux,
	}

	stdlog.Println(`starting server...`)

	// start the server
	err = server.ListenAndServe()
	if err != nil {
		stdlog.Printf("%s: error: %v", Environment, err)
	}

	stdlog.Println(`application exiting...`)
}
