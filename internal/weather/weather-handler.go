package weather

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"weatherserver/internal/model"
)

type API struct {
	weatherSvc WeatherService
}

func NewAPI(weatherSvc WeatherService) http.Handler {
	return &API{
		weatherSvc: weatherSvc,
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
			log.Printf("error: %v", err)
		}
	}(r.Body.Close)

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var data model.WeatherRecord
	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		log.Printf("error: %v, byte value: %s", err, bodyBytes)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = api.weatherSvc.Save(data)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, err.Error(), http.StatusMisdirectedRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(nil)
}
