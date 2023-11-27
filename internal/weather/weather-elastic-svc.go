package weather

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"weatherserver/internal/common"
	"weatherserver/internal/model"
)

const (
	weatherIndex = `test-indoor-weather`
)

type ElasticImpl struct {
	host   string
	client common.Client
}

func NewElasticService(client common.Client, elasticHost string) WeatherService {
	return &ElasticImpl{
		client: client,
		host:   elasticHost,
	}
}

func (svc *ElasticImpl) Save(data model.WeatherRecord) error {

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("error: %v", err)
		return err
	}
	body := io.NopCloser(bytes.NewReader(jsonBytes))

	url := fmt.Sprintf("%s/%s/_doc", svc.host, weatherIndex)
	log.Printf("calling: %s", url)
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		log.Printf("error: %v", err)
		return err
	}
	defer func(closeFunc func() error) {
		err = closeFunc()
		if err != nil {
			log.Printf("error: %v", err)
		}
	}(req.Body.Close)

	req.Header.Set(`Content-Type`, `application/json`)

	res, err := svc.client.Do(req)
	if err != nil {
		log.Printf("error: %v", err)
		return err
	}

	if res.StatusCode != http.StatusCreated {
		log.Printf("error: expected 201 status, got %d", res.StatusCode)
	}

	return err

}
