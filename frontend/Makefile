.PHONY: lint
lint: node_modules
	npm run lint

.PHONY: test
test: node_modules
	npm run compile
	CI=true npm test

node_modules: package.json
	npm install
	touch -m $@