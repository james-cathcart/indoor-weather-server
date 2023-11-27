package system

import (
	"encoding/json"
	"fmt"
	"net/http"
	"weatherserver/internal/logging"
)

type SystemHandler struct {
	log logging.Logger
}

func NewSystemHandler(logger logging.Logger) http.Handler {
	return &SystemHandler{
		log: logger,
	}
}

func (api *SystemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		api.handlGet(w, r)
	default:
		msg := `method not allowed`
		api.log.Info(msg)
		http.Error(w, msg, http.StatusMethodNotAllowed)
	}

}

func (api *SystemHandler) handlGet(w http.ResponseWriter, r *http.Request) {

	option := r.URL.Query().Get(`option`)
	switch option {
	case `health`:
		api.healthCheck(w, r)
	default:
		msg := fmt.Sprintf("invalid request option: %s", option)
		api.log.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
	}
}

func (api *SystemHandler) healthCheck(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, `method not allowed`, http.StatusMethodNotAllowed)
		return
	}

	data := struct {
		Message string `json:"message"`
	}{
		Message: `weather server is up`,
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		api.log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(`Content-Type`, `application/json`)
	byteCount, err := w.Write(jsonBytes)
	if err != nil {
		msg := fmt.Sprintf("wrote: %d bytes, `%s`", byteCount, err.Error())
		api.log.Error(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
}
