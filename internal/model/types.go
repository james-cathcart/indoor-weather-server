package model

type WeatherRecord struct {
	Timestamp   string  `json:"timestamp"`
	Humidity    int     `json:"humidity"`
	Temperature float64 `json:"temperature"`
	Pressure    float64 `json:"pressure"`
	Location    string  `json:"location"`
}
