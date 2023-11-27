FROM golang:1.21.0

ENV GOCACHE=/tmp/

ENV BASE_PATH=/go/src/github.com/james-cathcart/weather-server

WORKDIR $BASE_PATH

ADD . .