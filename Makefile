fix-mod:
	go mod tidy

run-server:
	go run cmd/server/main.go

run-client:
	go run cmd/client/main.go

up:
	docker compose up -d

down:
	docker compose down

migrate:
	go run cmd/migration/main.go

create-migration:
	migrate create -ext sql -dir migrations/ -seq change_this_text

build-client:

build-server:

gen-proto:
	protoc --go_out=internal/grpc/auth/proto  \
    		--go-grpc_out=internal/grpc/auth/proto \
    		internal/grpc/auth/proto/authService.proto

linter:
	go vet -vettool=./cmd/linter/linter ./...