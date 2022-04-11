TARGET=wolverine
VERSION=0.1.0

build: format test
	docker build --rm -t ${TARGET}:${VERSION} --build-arg VERSION=${VERSION} .

build-local: format test
	go build -o ${TARGET} cmd/main.go

format:
	gofmt -w .

clean:
	rm ${TARGET}

test:
	go test ./...

load-image:
	kind load docker-image --name wolverine ${TARGET}:${VERSION}

deploy:
	kubectl apply -f k8s/wolverine.yaml

.PHONY: format build build-local clean test load-image deploy
