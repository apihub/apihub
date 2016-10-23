.PHONY: test

DOCKERFILE_TEST := Dockerfile_test
IMAGE_TEST := apihub_test
PROJECT := github.com/apihub/apihub

all:
		go build -o apihub-api ./api/

###### Help ###################################################################

help:
		@echo '    all ................................. builds apihub'
		@echo '    deps ................................ installs dependencies'
		@echo '    docker-test ......................... runs tests in a container'
		@echo '    go-generate ......................... runs go generate'
		@echo '    go-vet .............................. runs go vet'
		@echo '    test ................................ runs tests locally'

###############################################################################

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

test: go-vet
	ginkgo -r -p -race -keepGoing .
