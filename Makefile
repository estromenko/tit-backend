all: dev

.PHONY: build
.SILENT: build
build: build-createsuperuser
	go build -o .tmp/main cmd/tit-backend/main.go

.PHONY: build-createsuperuser
.SILENT: build-createsuperuser
build-createsuperuser:
	go build -o ./.tmp/createsuperuser cmd/createsuperuser/main.go

.PHONY: createsuperuser
.SILENT: createsuperuser
createsuperuser: build-createsuperuser
	eval $$(cat .env) ./.tmp/createsuperuser

.PHONY: dev
.SILENT: dev
dev:
	eval $$(cat .env) $$(go env GOBIN)/air

.PHONY: format
.SILENT: format
format:
	$$(go env GOBIN)/gofumpt -l -w .

.PHONY: install-dev-requirements
install-dev-requirements:
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@master
	go install mvdan.cc/gofumpt@latest

.PHONY: install-migrate
install-migrate:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

.PHONY: lint
.SILENT: lint
lint:
	$$(go env GOBIN)/golangci-lint run

.PHONY: lint-fix
.SILENT: lint-fix
lint-fix:
	$$(go env GOBIN)/golangci-lint run --fix

.PHONY: start
.SILENT: start
start:
	eval $$(cat .env) ./.tmp/main

.PHONY: test
.SILENT: test
test:
	eval $$(cat .env) go test -v -cover ./... | grep -v 'no test files'
