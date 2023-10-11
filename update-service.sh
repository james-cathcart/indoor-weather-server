#!/bin/bash

echo "stopping service..."
sudo systemctl stop weather-server
echo "updating application..."
git pull
make
echo "starting service..."
sudo systemctl start weather-server
echo "update complete"