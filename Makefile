TARGET=wolverine

build: format
	go build -o ${TARGET} cmd/main.go

format:
	gofmt -w .

clean:
	rm ${TARGET}

.PHONY: format build clean
