.PHONY: dev prod build push deploy install test lint test_api lint_api test_frontend lint_frontend

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

test: test_api test_frontend

lint: lint_api lint_frontend

test_api:
	cd api && go test ./...

lint_api:
	cd api && golangci-lint run

test_frontend: install
	cd frontend && npm run compile && CI=true npm test

lint_frontend: install
	cd frontend && npm run lint
