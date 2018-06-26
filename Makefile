SHELL := /bin/bash -x
TRIGRD_APP := trigd
TRIGRT_APP := trigr
VERSION := `cat VERSION`

NATIVE_PLUGINS := example
# Support binary builds
PLATFORMS := linux darwin freebsd

all: build

.PHONY: clean
clean:
	rm -fr build 
	echo $(PLATFORMS)
	@- $(foreach PLAT,$(PLATFORMS), \
		mkdir -p build/$(PLAT) \
		)

deps:
	# Someday switch to vgo once it works with the code
	# go get -u golang.org/x/vgo

	# vgo build

	go get -u github.com/golang/dep/cmd/dep

vendor: deps
	dep ensure

.PHONY: build
build: clean 
	for app in $(TRIGRD_APP) $(TRIGRT_APP); do \
		cd cmd/$$app ; \
		for plat in $(PLATFORMS); do \
			echo "Building '$$app' for platform '$$plat' ..." ; \
			GOARCH=amd64 GOOS=$$plat go build -ldflags "-s -w" -o ../../build/$$plat/$$app ; \
		done; \
		cd ../../; \
	done

.PHONY: native-plugins
native-plugins: 
	for p in $(NATIVE_PLUGINS); do \
		echo $$p; \
		mkdir -p build/plugins/native; \
		cd native/example; go build -buildmode=plugin -o ../../build/plugins/native/example.so; \
	done

.PHONY: test
test: 
	go test ./...

.PHONY: test-integration
test-integration: 
	go test ./... --tags=integration


