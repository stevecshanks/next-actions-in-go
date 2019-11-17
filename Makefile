.PHONY: dev prod install test 

dev:
	docker-compose up

prod:
	docker-compose -f docker-compose-production.yml -p next-actions-prod up

install: frontend/node_modules

frontend/node_modules: frontend/package.json
	cd frontend && npm install && touch -m node_modules

test: go_tests frontend_tests

go_tests:
	cd api && go test

frontend_tests: install
	cd frontend && CI=true npm test
