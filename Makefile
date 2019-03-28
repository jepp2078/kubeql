export GO111MODULE=on
export GOOS:=$(shell go env GOOS)
export GOARCH:=$(shell go env GOARCH)
export BINNAME:=kubeql

all: build run

run:
	./bin/${BINNAME} --v=5 --logtostderr=true --kubeconfig=/home/${USER}/.kube/config

test:
	go test ./...

build-docker:
	docker run -it -e GOOS=${GOOS} -e GOARCH=${GOARCH} -v $(shell pwd):/${BINNAME} -w /${BINNAME} golang:1.12 make

build:
	go build -mod vendor -o bin/${BINNAME} .

generate: 
	go run github.com/99designs/gqlgen

clean:
	sudo rm bin/${BINNAME}
	