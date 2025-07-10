.PHONY: fmt
fmt:
	go mod tidy
	go fmt ./...
	go test ./...
	golangci-lint fmt
	golangci-lint run --fix