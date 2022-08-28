#!/bin/sh

docker build . -t abidh/postgres:latest

echo "**************"

docker images | grep abidh/postgres

docker-compose up