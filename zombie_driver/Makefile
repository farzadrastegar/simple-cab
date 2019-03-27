.PHONY: all test

all:
	go run cmd/main.go

build:
	go build cmd/main.go

test:
	make -C ./client test
	make -C ./http test
	make -C ./cmd test
	make -C ./config test

build-linux:
	export GOOS=linux
	go build -o api-server-linux-amd64
	export GOOS=darwin
