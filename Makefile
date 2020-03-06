all: vendor build run

vendor:
	go mod tidy
	go mod vendor

build:
	go build ./cmd/proxy

run:
	go run ./cmd/proxy

caller:
	grpcurl -v --plaintext -d '{"name":"Connor"}' localhost:8080 helloproto.Greeter/SayHello
	grpcurl --plaintext localhost:8080 healthproto.Health/Check

lint: ## Reorders imports and runs the golangci-lint checker
	goimports -d -e -w ./cmd ./pkg
	golangci-lint run

gen:
	protoc -I . -I external/ --go_out=plugins=grpc:. pkg/health/healthproto/health.proto
	protoc -I . -I external/ --grpc-gateway_out=logtostderr=true:. pkg/health/healthproto/health.proto
	protoc -I . --go_out=plugins=grpc:. pkg/hello/helloproto/helloworld.proto

.PHONY: all vendor build run caller gen

