COMMIT_SHA=$(shell git rev-parse HEAD)


all: vendor build run

vendor:
	go mod tidy
	go mod vendor

build:
	go build -ldflags="-X main.version=${COMMIT_SHA}" ./cmd/proxy

run:
	go run -ldflags="-X main.version=${COMMIT_SHA}" ./cmd/proxy

caller:
	grpcurl -v --plaintext -d '{"name":"Connor"}' localhost:8080 helloworld.Greeter/SayHello
	grpcurl --plaintext localhost:8080 grpc_health_proxy.Health/Check
	curl localhost:8082/healthz

lint: ## Reorders imports and runs the golangci-lint checker
	goimports -d -e -w ./cmd ./pkg/health/instrument.go ./pkg/health/server.go ./pkg/hello/server.go ./pkg/tools
	golangci-lint run

gen:
	protoc -I . -I external/ --go_out=plugins=grpc:pkg ./proto/health.proto
	protoc -I . -I external/ --go_out=plugins=grpc:pkg --grpc-gateway_out=logtostderr=true:pkg ./proto/health.proto
	protoc -I . --go_out=plugins=grpc:pkg ./proto/helloworld.proto
	mv pkg/proto/health.pb.go pkg/health/grpc_health_proxy/
	mv pkg/proto/health.pb.gw.go pkg/health/grpc_health_proxy/
	mv pkg/proto/helloworld.pb.go pkg/hello
	rmdir pkg/proto


.PHONY: all vendor build run caller gen

