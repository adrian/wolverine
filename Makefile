TARGET=wolverine

build: format
	go build -o ${TARGET} cmd/main.go

format:
	gofmt -w .

clean:
	rm ${TARGET}

test:
	go test ./...

.PHONY: format build clean test
