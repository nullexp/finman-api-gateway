.PHONY: dev-run install buf lint

export

install:

	@go mod tidy
	@go install github.com/bufbuild/buf/cmd/buf@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

buf:
	mkdir -p "./proto/user/v1"
	mkdir -p "./proto/transaction/v1"
	mkdir -p "./proto/auth/v1"
	curl -o ./proto/user/v1/user.proto https://raw.githubusercontent.com/nullexp/finman-user-service/main/proto/user/v1/user.proto
	curl -o ./proto/transaction/v1/transaction.proto https://raw.githubusercontent.com/nullexp/finman-transaction-service/main/proto/transaction/v1/transaction.proto
	@env PATH="$$PATH:$$(go env GOPATH)/bin" buf generate --template proto/buf.gen.yaml proto
	@echo "✅ buf done!"
	rm -rf "./proto/user"
	rm -rf "./proto/transaction"


buf-win:
	mkdir ".\proto\user\v1"
	mkdir ".\proto\transaction\v1"
	mkdir ".\proto\auth\v1"
	curl -o .\proto\user\v1\user.proto https://raw.githubusercontent.com/nullexp/finman-user-service/main/proto/user/v1/user.proto
	curl -o .\proto\transaction\v1\transaction.proto https://raw.githubusercontent.com/nullexp/finman-transaction-service/main/proto/transaction/v1/transaction.proto
	curl -o .\proto\auth\v1\auth.proto https://raw.githubusercontent.com/nullexp/finman-auth-service/main/proto/auth/v1/auth.proto
	@set PATH=%PATH%;%GOPATH%\bin
	@buf generate --template proto\buf.gen.yaml proto
	@echo "✅ buf done!"
	rmdir /S /Q ".\proto\user"
	rmdir /S /Q ".\proto\transaction"
	rmdir /S /Q ".\proto\auth"




run:
	go run ./cmd
	
lint:
	gofumpt -l -w .
	golangci-lint run  -v

test:
	go test ./...

docker-build:
	docker build -t finman-gateway-service .

docker-run:
	docker run -p 8081:8081 finman-gateway-service

docker-compose-up:
	docker-compose up --build 

docker-compose-down:
	docker-compose down --volumes