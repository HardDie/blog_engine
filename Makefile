.PHONY: default
default: help

.PHONY: help
help: ## help information about make commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: sqlc
sqlc: ## generate SQL methods
	sqlc generate

.PHONY: install
install: ## install dep
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/daixiang0/gci@latest

.PHONY: run-fe
run-fe: ## run dev server for fe
	docker-compose -f deployments/docker-compose.fe.yaml up -d

.PHONY: dependency
dependency: ## install dev dependency
	# gRPC generator
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

.PHONY: proto
proto: ## generate go files from *.proto
	protoc -I./pkg/proto/server \
		-I./pkg/proto \
		--go_out ./pkg/proto/server \
		--go_opt=paths=source_relative \
		--go-grpc_out ./pkg/proto/server \
		--go-grpc_opt=paths=source_relative \
		--grpc-gateway_out ./pkg/proto/server \
		--grpc-gateway_opt=paths=source_relative \
		--openapiv2_out ./ \
		--openapiv2_opt allow_merge=true,merge_file_name=api,omit_enum_default_value=true,output_format=yaml \
		./pkg/proto/server/*.proto


.PHONY: format
format: ## format code and imports
	go fmt ./...
	gci write -s standard -s default -s 'prefix(github.com/HardDie)' -s localmodule --skip-generated .
