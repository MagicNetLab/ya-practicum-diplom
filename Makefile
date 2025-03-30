fix-mod:
	go mod tidy

run-server:
	go run cmd/server/main.go

run-client:
	go run cmd/client/main.go run

up:
	docker compose up -d

down:
	docker compose down

migrate:
	go run cmd/migration/main.go

create-migration:
	migrate create -ext sql -dir migrations/ -seq change_this_text

# Определение переменных для разных платформ
BINARY_CLIENT=gophkeeper-client
BINARY_SERVER=gophkeeper-server
BUILD_DIR=bin
BUILD_DIR_CLIENT=$(BUILD_DIR)/client
BUILD_DIR_SERVER=$(BUILD_DIR)/server

# Создание директорий для сборки
build-dir:
	mkdir -p $(BUILD_DIR_CLIENT)
	mkdir -p $(BUILD_DIR_SERVER)

# Сборка клиента
build-client-linux: build-dir
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR_CLIENT)/$(BINARY_CLIENT)-linux-amd64 ./cmd/client/main.go

build-client-windows: build-dir
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR_CLIENT)/$(BINARY_CLIENT)-windows-amd64.exe ./cmd/client/main.go

build-client-macos: build-dir
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR_CLIENT)/$(BINARY_CLIENT)-darwin-amd64 ./cmd/client/main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR_CLIENT)/$(BINARY_CLIENT)-darwin-arm64 ./cmd/client/main.go

# Сборка сервера
build-server-linux: build-dir
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR_SERVER)/$(BINARY_SERVER)-linux-amd64 ./cmd/server/main.go

build-server-windows: build-dir
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR_SERVER)/$(BINARY_SERVER)-windows-amd64.exe ./cmd/server/main.go

build-server-macos: build-dir
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR_SERVER)/$(BINARY_SERVER)-darwin-amd64 ./cmd/server/main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR_SERVER)/$(BINARY_SERVER)-darwin-arm64 ./cmd/server/main.go

# Сборка для всех платформ
build-all-client: build-client-linux build-client-windows build-client-macos

build-all-server: build-server-linux build-server-windows build-server-macos

# Сборка всех приложений для всех платформ
build-all: build-all-client build-all-server

# Очистка директории сборки
clean-bin:
	rm -rf $(BUILD_DIR)

build-client: build-dir
	go build -o $(BUILD_DIR_CLIENT)/$(BINARY_CLIENT) ./cmd/client/main.go

build-server: build-dir
	go build -o $(BUILD_DIR_SERVER)/$(BINARY_SERVER) ./cmd/server/main.go

gen-proto:
	protoc --go_out=internal/grpc/auth/proto --go-grpc_out=internal/grpc/auth/proto internal/grpc/auth/proto/authService.proto
	protoc --go_out=internal/grpc/account/proto --go-grpc_out=internal/grpc/account/proto internal/grpc/account/proto/accountService.proto
	protoc --go_out=internal/grpc/card/proto --go-grpc_out=internal/grpc/card/proto internal/grpc/card/proto/cardService.proto
	protoc --go_out=internal/grpc/note/proto --go-grpc_out=internal/grpc/note/proto internal/grpc/note/proto/noteService.proto
	protoc --go_out=internal/grpc/files/proto --go-grpc_out=internal/grpc/files/proto internal/grpc/files/proto/fileService.proto

linter:
	go vet -vettool=./cmd/linter/linter ./...

mock:
	make mock-repo
	make mock-repo-models
	make mock-clients
	make mock-service

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
	mockery --name=CardSearchModel --dir=internal/repository/models --output=internal/repository/models/mocks --outpkg=mocks
	mockery --name=NoteModel --dir=internal/repository/models --output=internal/repository/models/mocks --outpkg=mocks
	mockery --name=NoteSearchModel --dir=internal/repository/models --output=internal/repository/models/mocks --outpkg=mocks
	mockery --name=FilesModel --dir=internal/repository/models --output=internal/repository/models/mocks --outpkg=mocks
	mockery --name=FilesSearchModel --dir=internal/repository/models --output=internal/repository/models/mocks --outpkg=mocks

mock-clients:
	mockery --name=S3Client --dir=internal/services/s3 --output=internal/services/s3/mocks --outpkg=mocks

mock-conf:
	mockery --name=ServerConfigurator --dir=internal/config --output=internal/config/mocks --outpkg=mocks
	mockery --name=SecretConfigurator --dir=internal/config --output=internal/config/mocks --outpkg=mocks
	mockery --name=S3Configurator --dir=internal/config --output=internal/config/mocks --outpkg=mocks
	mockery --name=JWTConfigurator --dir=internal/config --output=internal/config/mocks --outpkg=mocks
	mockery --name=DataBaseConfigurator --dir=internal/config --output=internal/config/mocks --outpkg=mocks
	mockery --name=AppConfigurator --dir=internal/config --output=internal/config/mocks --outpkg=mocks

mock-service:
	mockery --name=AppClient --dir=internal/services/client --output=internal/services/client/mocks --outpkg=mocks
	mockery --name=FileManager --dir=internal/services/client --output=internal/services/client/mocks --outpkg=mocks

mock-grpc:
	mockery --name=AuthServer --dir=internal/grpc/auth/proto --output=internal/grpc/auth/mocks --outpkg=mocks
	mockery --name=AuthClient --dir=internal/grpc/auth/proto --output=internal/grpc/auth/mocks --outpkg=mocks
	mockery --name=AccountsServer --dir=internal/grpc/account/proto --output=internal/grpc/account/mocks --outpkg=mocks
	mockery --name=AccountsClient --dir=internal/grpc/account/proto --output=internal/grpc/account/mocks --outpkg=mocks
	mockery --name=CardServer --dir=internal/grpc/card/proto --output=internal/grpc/card/mocks --outpkg=mocks
	mockery --name=CardClient --dir=internal/grpc/card/proto --output=internal/grpc/card/mocks --outpkg=mocks
	mockery --name=FilesServer --dir=internal/grpc/files/proto --output=internal/grpc/files/mocks --outpkg=mocks
	mockery --name=FilesClient --dir=internal/grpc/files/proto --output=internal/grpc/files/mocks --outpkg=mocks
	mockery --name=NoteServer --dir=internal/grpc/note/proto --output=internal/grpc/note/mocks --outpkg=mocks
	mockery --name=NoteClient --dir=internal/grpc/note/proto --output=internal/grpc/note/mocks --outpkg=mocks