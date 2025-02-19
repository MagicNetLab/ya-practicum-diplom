fix-mod:
	go mod tidy

run-server:
	go run cmd/server/main.go

run-client:
	go run cmd/client/main.go

migrate:

build-client:

build-server:

gen-proto:
	protoc --go_out=internal/grpc/auth/proto  \
    		--go-grpc_out=internal/grpc/auth/proto \
    		internal/grpc/auth/proto/authService.proto

linter: