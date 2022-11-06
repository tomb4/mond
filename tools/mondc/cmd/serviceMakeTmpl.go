package cmd

import "text/template"

var makefileTemplate, _ = template.New("").Parse(`
proto:
	protoc --go_out=./ --go_opt=paths=source_relative --go-grpc_out=./ --go-grpc_opt=paths=source_relative proto/{{.AppId}}.proto
	protoc --mond_out=./ --mond_opt=paths=source_relative  proto/{{.AppId}}.proto
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
