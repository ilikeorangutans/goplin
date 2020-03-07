.PHONY: all
all: run

.PHONY: run
run: build
	./goplin

.PHONY: build
build:
	go build -v ./cmd/goplin
