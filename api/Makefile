.PHONY: lint
lint: golangci-lint
	golangci-lint run

.PHONY: test
test:
	go test -race ./...

.PHONY: golangci-lint
golangci-lint: $(GOPATH)/bin/golangci-lint

$(GOPATH)/bin/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v1.27.0
