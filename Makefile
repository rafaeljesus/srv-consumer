.PHONY: all deps test build

all: deps test build

deps:
	@go get -u github.com/golang/dep/cmd/dep
	@dep ensure

test:
	@go vet ./{cmd,handler,platform,storage}/...
	@go test -v -race -cover ./{cmd,handler,platform,storage}/...

build:
	@GOBIN=/build go install -ldflags "-w -s" ./...
