all: clean

run:
	go run examples/download/main.go

test:
	go test -cover -v ./...

validate:
	@! gofmt -s -d -l . 2>&1 | grep -vE '^\.git/'
	go vet ./...

clean:
	go clean

.PHONY: build install test clean validate
