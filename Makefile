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
	protoc --go_out=internal/grpc/auth/proto --go-grpc_out=internal/grpc/auth/proto internal/grpc/auth/proto/authService.proto
	protoc --go_out=internal/grpc/account/proto --go-grpc_out=internal/grpc/account/proto internal/grpc/account/proto/accountService.proto

linter:
	go vet -vettool=./cmd/linter/linter ./...

mock-repo:
	mockery --name=AuthRepository --dir=internal/repository --output=internal/repository/mocks --outpkg=mocks
	mockery --name=AccountRepository --dir=internal/repository --output=internal/repository/mocks --outpkg=mocks

mock-repo-models:
	mockery --name=AccountModel --dir=internal/repository/models --output=internal/repository/models/mocks --outpkg=mocks
	mockery --name=AccountSearchModel --dir=internal/repository/models --output=internal/repository/models/mocks --outpkg=mocks
	mockery --name=TokenModel --dir=internal/repository/models --output=internal/repository/models/mocks --outpkg=mocks
	mockery --name=UserModel --dir=internal/repository/models --output=internal/repository/models/mocks --outpkg=mocks