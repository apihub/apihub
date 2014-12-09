help:
	@echo '    run-api .................. runs api server'
	@echo '    race ..................... runs race condition tests'
	@echo '    save-deps ................ generates the Godeps folder'
	@echo '    setup .................... sets up the environment'
	@echo '    test ..................... runs tests'

run-api:
	go run ./api/cmd/httpserver.go

race:
	go test $(GO_EXTRAFLAGS) -race -i ./...
	go test $(GO_EXTRAFLAGS) -race ./...

save-deps:
	$(GOPATH)/bin/godep save ./...

setup:
	go get $(GO_EXTRAFLAGS) -u -d -t ./...
	go get $(GO_EXTRAFLAGS) github.com/tools/godep
	$(GOPATH)/bin/godep restore ./...

test:
	go test ./...
