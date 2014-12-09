help:
	@echo '    race ..................... runs race condition tests'
	@echo '    test ..................... runs tests'
	@echo '    run-api .................. runs api server'

test:
	go test ./...

setup:
	go get $(GO_EXTRAFLAGS) -u -d -t ./...
	go get $(GO_EXTRAFLAGS) github.com/tools/godep
	$(GOPATH)/bin/godep restore ./...

save-deps:
	$(GOPATH)/bin/godep save ./...

run-api:
	go run ./api/cmd/httpserver.go

race:
	go test $(GO_EXTRAFLAGS) -race -i ./...
	go test $(GO_EXTRAFLAGS) -race ./...
