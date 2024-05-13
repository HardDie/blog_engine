.PHONY: default
default: help

.PHONY: help
help: ## help information about make commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: swagger
swagger: ## generate swagger for project
	swagger generate spec -m -o swagger.yaml

.PHONY: sqlc
sqlc: ## generate SQL methods
	sqlc generate

.PHONY: install
install: ## install dep
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

.PHONY: run-fe
run-fe: ## run dev server for fe
	docker-compose -f deployments/docker-compose.fe.yaml up -d
