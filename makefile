export SHELL:=/usr/bin/env bash -O extglob -c
export GO111MODULE:=on
export OS=$(shell uname | tr '[:upper:]' '[:lower:]')
export ARTIFACT=json-structure

build: GOOS ?= ${OS}
build: GOARCH ?= amd64
build: clean
build:
	GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags "-X main.buildTime=`date --iso-8601=s` -X main.buildVersion=`git rev-parse HEAD | cut -c-7`" .

test: clean
	go test -v -vet=all -failfast -race

clean:
	rm -f ${ARTIFACT}
	rm -f ${ARTIFACT}-*.txz
