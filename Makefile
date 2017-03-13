SHELL = /bin/bash

all: $(TARGETS)
	go build

install:
	go install

clean:
	go clean

shortr: src/*.go
	go get -v ./...
	go build -o shortr src/*.go
