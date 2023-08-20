tidy:
  go mod tidy
  go mod download

fmt: tidy
  go fmt ./...

test: tidy
  go test -race -v ./...

update:
  go get ./...

generate:
  go generate ./...
  wire

vulnerability-check:
  govulncheck ./...

clean-go:
  go clean