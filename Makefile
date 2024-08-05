MAKEFILE_PATH=$(shell readlink -f "${0}")
MAKEFILE_DIR=$(shell dirname "${MAKEFILE_PATH}")

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: fmt
fmt: tidy
	gofmt -s -w .
	goimports -local github.com/charming/charming -w .

.PHONY: lint
lint: fmt
	golangci-lint run ./...

.PHONY: test
test:
	go test -v -race -coverprofile=coverage.out ./...
	@echo "all tests passed"

.PHONY: test-report
test-report:
	@make test | grep -A 1 'coverage: '

.PHONY: build
build:
	CGO_ENABLED=0 go build -gcflags "all=-N -l" -o bin/charming cmd/charming/charming.go
