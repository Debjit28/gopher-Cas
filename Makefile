build:
	@go build -o bin/fs ./cmd/my-app

run: build
	@./bin/fs

test:
	@go test ./... -v