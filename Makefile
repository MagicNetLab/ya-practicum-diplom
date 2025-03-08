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
	protoc --go_out=internal/grpc/card/proto --go-grpc_out=internal/grpc/card/proto internal/grpc/card/proto/cardService.proto
	protoc --go_out=internal/grpc/note/proto --go-grpc_out=internal/grpc/note/proto internal/grpc/note/proto/noteService.proto
	protoc --go_out=internal/grpc/files/proto --go-grpc_out=internal/grpc/files/proto internal/grpc/files/proto/fileService.proto

linter:
	go vet -vettool=./cmd/linter/linter ./...

mock-repo:
	mockery --name=AuthRepository --dir=internal/repository --output=internal/repository/mocks --outpkg=mocks
	mockery --name=AccountRepository --dir=internal/repository --output=internal/repository/mocks --outpkg=mocks
	mockery --name=CardRepository --dir=internal/repository --output=internal/repository/mocks --outpkg=mocks
	mockery --name=NoteRepository --dir=internal/repository --output=internal/repository/mocks --outpkg=mocks
	mockery --name=FileRepository --dir=internal/repository --output=internal/repository/mocks --outpkg=mocks

mock-repo-models:
	mockery --name=AccountModel --dir=internal/repository/models --output=internal/repository/models/mocks --outpkg=mocks
	mockery --name=AccountSearchModel --dir=internal/repository/models --output=internal/repository/models/mocks --outpkg=mocks
	mockery --name=TokenModel --dir=internal/repository/models --output=internal/repository/models/mocks --outpkg=mocks
	mockery --name=UserModel --dir=internal/repository/models --output=internal/repository/models/mocks --outpkg=mocks
	mockery --name=CardModel --dir=internal/repository/models --output=internal/repository/models/mocks --outpkg=mocks
	mockery --name=NoteModel --dir=internal/repository/models --output=internal/repository/models/mocks --outpkg=mocks
	mockery --name=NoteSearchModel --dir=internal/repository/models --output=internal/repository/models/mocks --outpkg=mocks
	mockery --name=FilesModel --dir=internal/repository/models --output=internal/repository/models/mocks --outpkg=mocks
	mockery --name=FilesSearchModel --dir=internal/repository/models --output=internal/repository/models/mocks --outpkg=mocks

mock-clients:
	mockery --name=S3Client --dir=internal/services/s3 --output=internal/services/s3/mocks --outpkg=mocks