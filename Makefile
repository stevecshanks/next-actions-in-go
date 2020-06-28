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
	echo ${DOCKER_HUB_TOKEN} | ssh docker-deploy@${DOCKER_SERVER} "docker login -u stevecshanks --password-stdin"
	ssh docker-deploy@${DOCKER_SERVER} "docker-compose -f docker-compose-production.yml down --rmi all --remove-orphans || true"
	scp docker-compose-production.yml docker-deploy@${DOCKER_SERVER}:~/docker-compose-production.yml
	ssh docker-deploy@${DOCKER_SERVER} "source .env && docker-compose -f docker-compose-production.yml pull && docker-compose -f docker-compose-production.yml up --no-build -d"

.PHONY: test
test: $(SUBDIRS)

.PHONY: lint
lint: $(SUBDIRS)

.PHONY: $(SUBDIRS)
$(SUBDIRS):
	$(MAKE) -C $@ $(MAKECMDGOALS)