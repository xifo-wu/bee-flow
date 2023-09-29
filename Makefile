SHELL = bash

VERSION := $(shell cat VERSION)

build:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build --ldflags "-s -X bee-flow/cmd/version.VERSION=$(VERSION)" -o beeflow main.go
