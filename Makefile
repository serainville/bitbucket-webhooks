default: test coverage

test:
	go test -coverprofile=coverage.out ./...

coverage:
	go tool cover -func=coverage.out