default: lint vet test coverage

lint:
	golint
	golangci-lint run

vet:
	go vet

test:
	go test -coverprofile=coverage.out ./...

coverage:
	go tool cover -func=coverage.out