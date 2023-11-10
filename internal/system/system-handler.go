package system

import (
	"encoding/json"
	"net/http"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {

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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(`Content-Type`, `application/json`)
	w.Write(jsonBytes)
}
