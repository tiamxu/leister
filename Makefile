.PHONY: build run test lint clean

build:
	go build -o gigctl .

run: build
	./gigctl

test:
	go test ./...

clean:
	rm -f gigctl
