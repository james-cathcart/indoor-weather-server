FROM golang:1.21.0 AS builder

ENV GOCACHE=/tmp

ENV BASE_PATH=/go/src/github.com/james-cathcart/weather-server
ENV CODE_DIR=code

RUN mkdir /weather-server

WORKDIR $BASE_PATH

ADD . $CODE_DIR

WORKDIR $BASE_PATH/$CODE_DIR

RUN go build -o /bin/weather-server cmd/api/main.go

FROM ubuntu

RUN apt-get update && apt-get install -y gnupg2
RUN apt-get install wget -y
RUN wget -qO - https://artifacts.elastic.co/GPG-KEY-elasticsearch | apt-key add -
RUN apt-get install apt-transport-https
RUN echo "deb https://artifacts.elastic.co/packages/7.x/apt stable main" | tee -a /etc/apt/sources.list.d/elastic-7.x.list
RUN apt-get update && apt-get install filebeat -y
RUN service filebeat start

RUN apt-get install ca-certificates -y
COPY --from=builder /bin/weather-server .
COPY --chmod=644 filebeat/filebeat.yml /etc/filebeat/filebeat.yml
COPY startup.sh .


EXPOSE 8080

ENTRYPOINT ["./startup.sh"]