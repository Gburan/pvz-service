LINTER = golangci-lint
LINTER_FLAGS = run

.DEFAULT_GOAL := lint

.PHONY: lint
lint:
	$(LINTER) $(LINTER_FLAGS)

.PHONY: lint-fix
lint-fix:
	$(LINTER) $(LINTER_FLAGS) --fix

.PHONY: install
install-linter:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install go.uber.org/mock/mockgen@latest

.PHONY: run-test-clean
run-test-clean:
	go clean -testcache
	go test ./...

.PHONY: gen
gen:
	go generate ./...

.PHONY: run-compose-f
run-compose-f:
	go mod tidy
	go generate ./...
	docker-compose up -d

.PHONY: run-compose
run-compose:
	docker-compose up -d

.PHONY: run-compose-b
run-compose-b:
	docker-compose up -d --build
