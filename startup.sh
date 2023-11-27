#!/bin/bash

touch log.json
service filebeat restart

./weather-server