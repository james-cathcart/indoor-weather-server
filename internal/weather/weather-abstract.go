package weather

import "weatherserver/internal/model"

type WeatherService interface {
	Save(data model.WeatherRecord) error
}
