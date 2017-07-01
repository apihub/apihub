.PHONY: test

DOCKERFILE_TEST := Dockerfile_test
IMAGE_TEST := apihub_test
PROJECT := github.com/apihub/apihub

all:
		go build -o apihub-api ./cmd/api/
		go build -o apihub-gateway ./cmd/gateway/

###### Help ###############################################################################

help:
		@echo '    all ................................. builds apihub'
		@echo '    concourse ........................... runs concourse in a local docker container'
		@echo '    deps ................................ installs dependencies'
		@echo '    docker-test ......................... runs tests in a container'
		@echo '    go-generate ......................... runs go generate'
		@echo '    go-vet .............................. runs go vet'
		@echo '    pipeline ............................ setups pipeline'
		@echo '    setup ............................... sets up the dev environment'
		@echo '    stop-concourse ...................... stops local concourse'
		@echo '    test ................................ runs tests locally'

###### Help ###############################################################################

concourse:
	cd ci/setup \
	; docker-compose up

deps:
	glide up

docker-test: image_test
	docker run -t --privileged --rm -v $(PWD):/go/src/$(PROJECT) $(IMAGE_TEST) make test

image_test:
	docker build -t $(IMAGE_TEST) -f $(DOCKERFILE_TEST) .

go-generate:
	go generate `go list ./... | grep -v vendor`

go-vet:
	go vet `go list ./... | grep -v vendor`

pipeline:
	./ci/set-pipeline

setup: deps
	cd vendor/github.com/hashicorp/consul \
	; CONSUL_DEV=true make \
	; mv bin/consul $(GOPATH)/bin

stop-concourse:
	docker ps | grep concourse | awk '{print $1}' | xargs docker stop

test: go-vet
	ginkgo -r -p -race -keepGoing .
