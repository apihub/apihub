FROM golang:latest

ADD . /go/src/github.com/apihub/apihub
WORKDIR /go/src/github.com/apihub/apihub

RUN apt-get update && apt-get install -y supervisor
ADD supervisord.conf /etc/supervisor/conf.d/supervisord.conf

RUN go get github.com/tools/godep
RUN godep restore ./...
RUN cd example && go build -o /go/bin/apihub_api
RUN cd gateway/example && go build -o /go/bin/apihub_gateway
#ENTRYPOINT ["/go/bin/apihub"]
ENTRYPOINT ["/usr/bin/supervisord"]