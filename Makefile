SHELL = /bin/bash

TARGETS = shortr

all: $(TARGETS)

clean:
	go clean
	rm -f $(TARGETS)

$(TARGETS): src/*.go
	go get -v ./...
	go build -o $@ src/*.go
