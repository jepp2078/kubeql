export GO111MODULE=on
export GOOS:=$(shell go env GOOS)
export GOARCH:=$(shell go env GOARCH)
export BINNAME:=kubeql

all: build run

run:
	./bin/${BINNAME} --v=5 --logtostderr=true

test:
	go test ./...

build-docker:
	docker run -it -e GOOS=${GOOS} -e GOARCH=${GOARCH} -v $(shell pwd):/${BINNAME} -w /${BINNAME} golang:1.11 make build

build:
	go build -mod vendor -o bin/${BINNAME} .

generate: 
	go run github.com/99designs/gqlgen

clean:
	sudo rm bin/${BINNAME}
	