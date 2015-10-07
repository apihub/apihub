FROM golang:latest

ADD . /go/src/github.com/apihub/apihub
WORKDIR /go/src/github.com/apihub/apihub

RUN go get github.com/tools/godep
RUN godep restore ./...
RUN cd example && go build -o /go/bin/apihub
ENTRYPOINT ["/go/bin/apihub"]