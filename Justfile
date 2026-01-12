IMAGE := "kuberhealthy/http-check"
TAG := "latest"

# Build the http check container locally.
build:
	podman build -f Containerfile -t {{IMAGE}}:{{TAG}} .

# Run the unit tests for the http check.
test:
	go test ./...

# Build the http check binary locally.
binary:
	go build -o bin/http-check ./cmd/http-check
