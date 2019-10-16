export GOPATH := $(shell pwd)
export PATH := $(GOPATH)/bin:$(PATH)

buildes:
	go build -v -o bin/mysql-elasticsearch ./src/datatravel/elasticsearch.go

build: LDFLAGS   += $(shell GOPATH=${GOPATH} src/build/ldflags.sh)
build:
	@echo "--> Building..."
	@mkdir -p bin/
	go build -v -o bin/datatravel --ldflags '$(LDFLAGS)' src/datatravel/datatravel.go
	go build -v -o bin/datatravelcli --ldflags '$(LDFLAGS)' src/cli/cli.go
	@chmod 755 bin/*

clean:
	@echo "--> Cleaning..."
	@mkdir -p bin/
	@go clean
	@rm -f bin/datatravel bin/datatravelcli bin/mysql-elasticsearch

fmt:
	go fmt ./...

test:
	@echo "--> Testing..."
	@$(MAKE) testshift
	@$(MAKE) testcanal
	@$(MAKE) testelasticsearch

testshift:
	go test -v -race shift

testcanal:
	go test -v vendor/github.com/siddontang/go-mysql/canal/...
	go test -v vendor/github.com/siddontang/go-mysql/replication/...

testelasticsearch:
	go test -v -race elasticsearch/elastic
	go test -v -race elasticsearch/river

# code coverage
COVPKGS =	shift\
			vendor/github.com/siddontang/go-mysql/canal/...
coverage:
	go build -v -o bin/gotestcover \
	src/vendor/github.com/pierrre/gotestcover/*.go;
	gotestcover -coverprofile=coverage.out -v $(COVPKGS)
	go tool cover -html=coverage.out
.PHONY: build clean install fmt test coverage .go-version
