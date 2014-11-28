help:
	@echo '    race ..................... runs race condition tests'
	@echo '    test ..................... runs tests'

test:
	go test ./...

race:
	go test $(GO_EXTRAFLAGS) -race -i ./...
	go test $(GO_EXTRAFLAGS) -race ./...
