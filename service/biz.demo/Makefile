
proto:
	protoc -Iproto/ -I./../../../ --go_out=plugins=grpc:proto proto/Bizdemo.proto
	protoc -Iproto/ -I./../../../ --meta_out=plugins=meta:proto proto/Bizdemo.proto
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

