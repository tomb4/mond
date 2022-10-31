package cmd

import "text/template"

var makefileTemplate, _ = template.New("").Parse(`
proto:
	protoc -Iproto/ -I./../../../ --go_out=plugins=grpc:proto proto/{{.AppId}}.proto
	protoc -Iproto/ -I./../../../ --meta_out=plugins=meta:proto proto/{{.AppId}}.proto
run:
	go run cmd/server.go
server:
	go build -o 'bin/server'  ./cmd/server.go

build:
	GOOS=linux GOARCH=amd64 go build -o 'bin/server' ./cmd/server.go

develop:

test:
	../../scripts/build.sh {{.App1Id}} {{.Port}} test

.PHONY: server
.PHONY: proto

`)
