 SHELL := /bin/bash

# This would be a nice discover plugin that would traverse the code base and vendor external packages
GO_DEPS = github.com/GeertJohan/go.rice \
					github.com/GeertJohan/go.rice/rice \
					bitbucket.org/lateefj/httphacks \
					github.com/Sirupsen/logrus \
					github.com/lateefj/trigr \
					github.com/lateefj/trigr/ext \
					golang.org/x/net/websocket

GB = $(GOPATH)/bin/gb

all: build ui-build embed

clean:
	rm bin/*

ui-deps:
	npm install -g bower
	npm install -g vulcanize # TODO figure out a way to get polymer to be packagable

ui-build:
	cd src/cmd/gopolyd/ui/ && \
		bower install ; \
		ls bower_components

deps:
	go get -u github.com/constabulary/gb/...
	go get -u github.com/GeertJohan/go.rice 

vendor-update:
	for dep in $(GO_DEPS) ; do \
		echo $$dep ; \
		$(GB) vendor update $$dep ; \
	done

vendor: deps
	for dep in $(GO_DEPS) ; do \
		echo $$dep ; \
		$(GB) vendor fetch $$dep ; \
	done

build: deps
	mkdir -p bin
	# Linux build
	GOARCH=amd64 GOOS=linux gb build
	# OS X build
	GOARCH=amd64 GOOS=darwin gb build 

embed: build ui-build
	cd src/cmd/gopolyd/ && rice append --exec ../../../bin/linux_gopolymerd
	cd src/cmd/gopolyd/ && rice append --exec ../../../bin/darwin_gopolymerd


# Crosscompile sadness :(
build-linux: deps
	cd src/cmd/gopolyd && env GOOS=linux go build -o ../../../bin/gopolyd

embed-linux: build-linux
	cd src/cmd/gopolyd/ && rice append --exec ../../../bin/gopolyd
	
