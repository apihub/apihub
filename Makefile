define HG_ERROR
FATAL: You need Mercurial (hg) to download backstage dependencies.
endef

define GIT_ERROR
FATAL: You need Git to download backstage dependencies.
endef

help:
	@echo '    doc ...................... generates a new doc version'
	@echo '    run-api .................. runs api server'
	@echo '    race ..................... runs race condition tests'
	@echo '    save-deps ................ generates the Godeps folder'
	@echo '    setup .................... sets up the environment'
	@echo '    test ..................... runs tests'

doc:
	@cd docs && make clean && make html SPHINXOPTS="-W"

run-api:
	go run ./api/cmd/httpserver.go

race:
	go test $(GO_EXTRAFLAGS) -race -i ./...
	go test $(GO_EXTRAFLAGS) -race ./...

save-deps:
	$(GOPATH)/bin/godep save ./...

setup:
	$(if $(shell hg), , $(error $(HG_ERROR)))
	$(if $(shell git), , $(error $(GIT_ERROR)))
	go get $(GO_EXTRAFLAGS) -u -d -t ./...
	go get $(GO_EXTRAFLAGS) github.com/tools/godep
	$(GOPATH)/bin/godep restore ./...

test:
	go test ./...