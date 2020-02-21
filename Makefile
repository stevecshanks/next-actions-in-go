.PHONY: dev prod install test build push deploy

dev:
	source .env && docker-compose up

prod: build
	docker-compose -f docker-compose-production.yml

build:
	docker-compose -f docker-compose-production.yml -f docker-compose-production-build.yml build

push:
	echo ${DOCKER_HUB_TOKEN} | docker login -u stevecshanks --password-stdin
	docker-compose -f docker-compose-production.yml -f docker-compose-production-build.yml push

deploy:
	echo ${DOCKER_HUB_TOKEN} | ssh docker-deploy@${DOCKER_SERVER} "docker login -u stevecshanks --password-stdin"
	ssh docker-deploy@${DOCKER_SERVER} "docker-compose -f docker-compose-production.yml down --rmi all --remove-orphans || true"
	scp docker-compose-production.yml docker-deploy@${DOCKER_SERVER}:~/docker-compose-production.yml
	ssh docker-deploy@${DOCKER_SERVER} "source .env && docker-compose -f docker-compose-production.yml pull && docker-compose -f docker-compose-production.yml up --no-build -d"

install: frontend/node_modules

frontend/node_modules: frontend/package.json
	cd frontend && npm install && touch -m node_modules

test: go_tests frontend_tests

go_tests:
	cd api && go test ./...

frontend_tests: install
	cd frontend && CI=true npm test
