package model

type WeatherRecord struct {
	Timestamp         string  `json:"_timestamp"`
	CardinalDirection string  `json:"cardinal_direction"`
	Humidity          int     `json:"humidity"`
	Temperature       float64 `json:"temperature"`
	TemperatureF      float64 `json:"temperature_f"`
	Pressure          float64 `json:"pressure"`
	Location          string  `json:"location"`
}
