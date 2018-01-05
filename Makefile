appname := promqtt

sources = $(shell find . -type f -name '*.go' -not -path "./vendor/*")
build_version = $(shell date +%Y%m%d-%H%M%S)+$(shell git rev-parse --short HEAD)

build = GOOS=$(1) GOARCH=$(2) go build -ldflags "-X=main.build=$(build_version)" -o build/$(appname)$(3)
tar = cd build && tar -cvzf $(appname).$(1)-$(2).tar.gz $(appname)$(3) && rm $(appname)$(3)
zip = cd build && zip $(appname).$(1)-$(2).zip $(appname)$(3) && rm $(appname)$(3)

.PHONY: all test clean fmt get windows darwin linux

all: get linux windows darwin

test:
	./test/run-integration-tests.sh

clean:
	rm -rf build/

fmt:
	@gofmt -l -w $(sources)

get:
	@echo ">> Fetching dependencies"
	go get -v

##### LINUX #####
linux: build/$(appname).linux-amd64.tar.gz

build/$(appname).linux-amd64.tar.gz: $(sources)
	$(call build,linux,amd64,)
	$(call tar,linux,amd64)

##### DARWIN (MAC) #####
darwin: build/$(appname).darwin-amd64.tar.gz

build/$(appname).darwin-amd64.tar.gz: $(sources)
	$(call build,darwin,amd64,)
	$(call tar,darwin,amd64)

##### WINDOWS #####
windows: build/$(appname).windows-amd64.zip

build/$(appname).windows-amd64.zip: $(sources)
	$(call build,windows,amd64,.exe)
	$(call zip,windows,amd64,.exe)
