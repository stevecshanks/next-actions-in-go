#!/bin/bash

set -e

source .env

docker login -u stevecshanks --password-stdin
docker-compose -f docker-compose-production.yml pull
docker-compose -f docker-compose-production.yml down --remove-orphans
docker-compose -f docker-compose-production.yml up --no-build -d
