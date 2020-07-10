SUBDIRS = api frontend

.PHONY: dev
dev:
	source .env && docker-compose up

.PHONY: prod
prod: build
	docker-compose -f docker-compose-production.yml

.PHONY: build
build:
	docker-compose -f docker-compose-production.yml -f docker-compose-production-build.yml build

.PHONY: push
push:
	echo ${DOCKER_HUB_TOKEN} | docker login -u stevecshanks --password-stdin
	docker-compose -f docker-compose-production.yml -f docker-compose-production-build.yml push

.PHONY: deploy
deploy:
	scp docker-compose-production.yml docker-deploy.sh docker-deploy@${DOCKER_SERVER}:~/
	echo ${DOCKER_HUB_TOKEN} | ssh docker-deploy@${DOCKER_SERVER} "/bin/bash docker-deploy.sh"

.PHONY: test
test: $(SUBDIRS)

.PHONY: lint
lint: $(SUBDIRS)

.PHONY: $(SUBDIRS)
$(SUBDIRS):
	$(MAKE) -C $@ $(MAKECMDGOALS)