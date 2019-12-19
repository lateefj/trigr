SHELL := /bin/bash -x
TRIGRD_APP := trigd
TRIGRT_APP := trigr
VERSION := `cat VERSION`

# Support binary builds
PLATFORMS := linux darwin freebsd

all: build

.PHONY: clean
clean:
	rm -rf build 
	echo $(PLATFORMS)
	@- $(foreach PLAT,$(PLATFORMS), \
		mkdir -p build/$(PLAT); \
		)


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

.PHONY: test
test: 
	go test ./...

.PHONY: test-integration
test-integration: 
	go test ./... --tags=integration


