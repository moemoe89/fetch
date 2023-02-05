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
	go test -v -race -cover -coverprofile=.coverage/pkg.coverage.tmp ./pkg/... && \
	cat .coverage/pkg.coverage.tmp | grep -v "_mock.go\|_mockgen.go" > .coverage/pkg.coverage && \
	go tool cover -html=.coverage/pkg.coverage -o .coverage/html/pkg.coverage.html;
	rm .coverage/pkg.coverage .coverage/pkg.coverage.tmp

clean:
	go clean
