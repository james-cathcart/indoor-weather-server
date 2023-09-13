all: clean build

build:
	go build -o bin/weather-server cmd/api/main.go

clean:
	rm -f bin/weather-server

run:
	go run cmd/api/main.go