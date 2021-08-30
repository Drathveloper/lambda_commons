build:
	go build ./...

generate-mocks:
	go generate ./...

fmt:
	go fmt ./...

test-only:
	go test ./...

test:
	generate-mocks
	test-only