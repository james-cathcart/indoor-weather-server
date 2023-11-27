package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"weatherserver/internal/logging"
	"weatherserver/internal/model"
)

type API struct {
	weatherSvc WeatherService
	log        logging.Logger
}

func NewAPI(weatherSvc WeatherService, logger logging.Logger) http.Handler {
	return &API{
		weatherSvc: weatherSvc,
		log:        logger,
	}
}

func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		api.handlePost(w, r)
	default:
		http.Error(w, `method not allowed`, http.StatusMethodNotAllowed)
	}
}

func (api *API) handlePost(w http.ResponseWriter, r *http.Request) {

	defer func(closeFunc func() error) {
		err := closeFunc()
		if err != nil {
			api.log.Error(err.Error())
		}
	}(r.Body.Close)

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		api.log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var data model.WeatherRecord
	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		api.log.Error(fmt.Sprintf("error: %v, byte value: %s", err, bodyBytes))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = api.weatherSvc.Save(data)
	if err != nil {
		api.log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusMisdirectedRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	written, err := w.Write(nil)
	if err != nil {
		api.log.Error(fmt.Sprintf("error: %d bytes written, message: `%v`", err, written))
		return
	}
}
