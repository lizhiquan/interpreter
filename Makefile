test:
	@echo "==> Running tests..."
	@go clean -testcache ./...
	@go test `go list ./... | grep -v cmd` -race -p 1 --cover
