package weather

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"weatherserver/internal/common"
	"weatherserver/internal/logging"
	"weatherserver/internal/model"
)

var (
	ElasticIndex string
)

type ElasticImpl struct {
	host   string
	client common.Client
	log    logging.Logger
}

func NewElasticService(client common.Client, elasticHost string, logger logging.Logger) WeatherService {
	return &ElasticImpl{
		client: client,
		host:   elasticHost,
		log:    logger,
	}
}

func (svc *ElasticImpl) Save(data model.WeatherRecord) error {

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		svc.log.Error(err.Error())
		return err
	}
	body := io.NopCloser(bytes.NewReader(jsonBytes))

	url := fmt.Sprintf("%s/%s/_doc", svc.host, ElasticIndex)
	svc.log.Info(fmt.Sprintf("calling: %s", url))
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		svc.log.Error(err.Error())
		return err
	}
	defer func(closeFunc func() error) {
		err = closeFunc()
		if err != nil {
			svc.log.Error(err.Error())
		}
	}(req.Body.Close)

	req.Header.Set(`Content-Type`, `application/json`)

	res, err := svc.client.Do(req)
	if err != nil {
		svc.log.Error(err.Error())
		return err
	}

	if res.StatusCode != http.StatusCreated {
		svc.log.Error(fmt.Sprintf("expected 201 status, got %d", res.StatusCode))
	}

	return err

}
