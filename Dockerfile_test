FROM albertoleal/wily
MAINTAINER Alberto Leal <albertonb@gmail.com>

RUN set -x \
  && go get github.com/onsi/ginkgo/ginkgo \
  && go get github.com/onsi/gomega

RUN set -x \
  && git clone https://github.com/hashicorp/consul /go/src/github.com/hashicorp/consul \
	&& cd /go/src/github.com/hashicorp/consul \
	&& CONSUL_DEV=true make \
	&& mv bin/consul /go/bin

WORKDIR /go/src/github.com/apihub/apihub/
