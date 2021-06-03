.PHONY: all clean build test test-unit

all: build

clean:
	rm -rf bin/

test: test-unit

test-unit:
	go test -v -race -tags=unit -coverprofile=coverage.txt -covermode=atomic ./...

build:
	docker build -t scratch-service .

run:
	docker run -p 8000:8000 -it scratch-service