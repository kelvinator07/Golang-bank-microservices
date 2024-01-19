#!/bin/sh

echo "start"

docker compose down && docker rmi golang-bank-microservices-api 2> /dev/null || true && docker compose up 
