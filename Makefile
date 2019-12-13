.PHONY: dev prod install test build push deploy

dev:
	docker-compose up

prod:
	docker-compose -f docker-compose-production.yml -p next-actions-prod up

build:
	docker-compose -f docker-compose-production.yml build

push:
	echo ${DOCKER_HUB_TOKEN} | docker login -u stevecshanks --password-stdin && docker-compose -f docker-compose-production.yml push

deploy: build push

install: frontend/node_modules

frontend/node_modules: frontend/package.json
	cd frontend && npm install && touch -m node_modules

test: go_tests frontend_tests

go_tests:
	cd api && go test

frontend_tests: install
	cd frontend && CI=true npm test
