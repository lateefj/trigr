 SHELL := /bin/bash


all: build ui-build embed

.PHONY: clean
clean:
	rm bin/*

.PHONY: test
test:
	go test ./...

build: 
	mkdir -p bin
	# Linux build
	GOARCH=amd64 GOOS=linux gb build
	# OS X build
	GOARCH=amd64 GOOS=darwin gb build 



