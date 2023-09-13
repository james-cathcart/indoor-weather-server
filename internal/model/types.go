package model

type WeatherRecord struct {
	TimeStamp   int64   `json:"time"`
	Humidity    int     `json:"humidity"`
	Temperature float64 `json:"temperature"`
	Pressure    float64 `json:"pressure"`
	Location    string  `json:"location"`
}
