# Dev Targets
GO_PROJECT_PATH = /go/src/github.com/jfbramlett/go-aop
DOCKER_TAG = orgunitsmigrator

# Default Target

.PHONY: compile
compile: lint vendor
	go build -mod=vendor -o bin/orgunitsmigrator ./pkg/...

# Tools
.PHONY: lint
lint:
	docker run --rm -t --entrypoint=linter -v `pwd`:$(GO_PROJECT_PATH) -w $(GO_PROJECT_PATH) registry.namely.land/namely/golang:dev-latest

vendor:
	go mod vendor

.PHONY: test
test: vendor
	go test -cover ./pkg/...
