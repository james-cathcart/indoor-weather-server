# Indoor Weather Server

This repository supports the Raspberry Pi SenseHat Indoor Weather application

## Overview
This is a simple application designed to collect weather data from _n_ number of nodes. The data will be persisted to Elasticsearch and utilize Kibana to visualize the data. This is meant partially as a tool to understand the fluxuations of temperature and humidity throughout an indoor environment, but also as a way to experiment with Kibana and other technologies requiring a dataset.

# Setup

## Configuration
You will need to following environment variables:
- **_ELASTIC_HOST_** - full hostname/port for the Elasticsearch host (ex: `http://elastic:9200`)
## Installation
Just clone the repository
## Run via CLI
```
make && bin/weather-server
```

## Run as Service
Add the following to `/lib/systemd/system/weather-server.service` updating the path to the server binary for the `ExecStart` value.
```
[Unit]
Description=Weather Server for Indoor Weather application

[Service]
Type=simple
Environment="ELASTIC_HOST=http://ELASTIC_HOST:ELASTIC_PORT"
ExecStart=/path/to/weather-server

[Install]
WantedBy=multi-user.target
```
Start the service and enable start-on-boot

# Development
## Architecture
![](docs/classes.png)

## Deployment
![](docs/deployment.png)

## Logic & Workflow
### Node/Server Flow
![](docs/server-flow.png)

### User Flow
![](docs/user-flow.png)

## Example Elasticsearch Response
```
        ...
        "hits": [
            {
                "_index": "test-indoor-weather",
                "_type": "_doc",
                "_id": "tvP2jooBkamCBFXjrHTK",
                "_score": 1.0,
                "_source": {
                    "time": 1694615645,
                    "humidity": 45,
                    "temperature": 24.984375,
                    "pressure": 1012.997802734375,
                    "location": "office"
                }
            },
            {
                "_index": "test-indoor-weather",
                "_type": "_doc",
                "_id": "t_P2jooBkamCBFXjtHSv",
                "_score": 1.0,
                "_source": {
                    "time": 1694615647,
                    "humidity": 45,
                    "temperature": 24.984375,
                    "pressure": 1013.00634765625,
                    "location": "office"
                }
            },
            ...
```