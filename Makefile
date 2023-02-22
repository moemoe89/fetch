GOPATH=$(shell go env GOPATH)

GOLANGCI_LINT_VERSION=v1.50.1

install-gomock:
	@echo "\n>>> Install gomock\n"
	go install github.com/golang/mock/mockgen

install-linter:
	@echo "\n>>> Install GolangCI-Lint"
	@echo ">>> https://github.com/golangci/golangci-lint/releases \n"
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/${GOLANGCI_LINT_VERSION}/install.sh | \
	sh -s -- -b ${GOPATH}/bin ${GOLANGCI_LINT_VERSION}

lint:
	@echo "\n>>> Run GolangCI-Lint\n"
	/bin/bash ./scripts/lint.sh

test:
	mkdir -p .coverage/html
	go test -v -race -cover -coverprofile=.coverage/pkg.coverage ./pkg/... && \
	cat .coverage/pkg.coverage | grep -v "_mock.go\|_mockgen.go" > .coverage/pkg.mockless.coverage && \
	go tool cover -html=.coverage/pkg.mockless.coverage -o .coverage/html/pkg.coverage.html;

mock:
	@echo "\n>>> Generates Mock\n"
	go generate ./...

.PHONY: build

build:
	go build -o fetch ./cmd

docker-build:
	docker build -t fetch -f ./build/Dockerfile .

docker-run:
	docker run --rm -it fetch sh

clean:
	go clean
