GO_VERSION := 1.22

TAG := $(shell git describe --abbrev=0 --tags --always)
HASH := $(shell git rev-parse HEAD)
DATE := $(shell date +%Y-%m-%dT%H:%M:%S)

LDFLAGS := -w -X github.com/kazakh-in-nz/hello-api/handlers.hash=$(HASH) -X github.com/kazakh-in-nz/hello-api/handlers.tag=$(TAG) -X github.com/kazakh-in-nz/hello-api/handlers.date=$(DATE)

.PHONY: install-go init-go

setup: install-go init-go install-lint copy-hooks install-godog

install-godog:
	go install github.com/cucumber/godog/cmd/godog@latest

install-go:
	wget "https://golang.org/dl/go$(GO_VERSION).linux-amd64.tar.gz"
	sudo tar -C /usr/local -xzf go$(GO_VERSION).linux-amd64.tar.gz
	rm go$(GO_VERSION).linux-amd64.tar.gz

init-go:
	echo 'export PATH=$$PATH:/usr/local/go/bin' >> $${HOME}/.zshrc
	echo 'export PATH=$$PATH:$${HOME}/go/bin' >> $${HOME}/.zshrc

build:
	go build -ldflags "$(LDFLAGS)" -o api cmd/main.go

test:
	go test -tags="unit" -v ./... -coverprofile=coverage.out

test-bdd:
	cd ./cmd/features && go test -tags=bdd -v -godog.paths ./*.feature -run "^TestFeatures/"

coverage:
	go tool cover -func coverage.out | grep "total:" | \
	awk '{print ((int($$3) > 80) != 1)}'

report:
	go tool cover -html=coverage.out -o cover.html

check-format:
	test -z $$(go fmt ./...)

vet:
	test -z $$(go vet ./...)

install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint:
	golangci-lint run

copy-hooks:
	chmod +x scripts/hooks/*
	cp -r scripts/hooks/* .git/.
