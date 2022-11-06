
proto:
	protoc --go_out=./ --go_opt=paths=source_relative --go-grpc_out=./ --go-grpc_opt=paths=source_relative proto/Bizdemo.proto
	protoc --mond_out=./ --mond_opt=paths=source_relative  proto/Bizdemo.proto
run:
	go run cmd/server.go

server:
	go build -o 'bin/server'  ./cmd/server.go

build:
	GOOS=linux GOARCH=amd64 go build -o 'bin/server' ./cmd/server.go

develop:

test:
	../../scripts/build.sh biz-demo 19001 test

.PHONY: server
.PHONY: proto

