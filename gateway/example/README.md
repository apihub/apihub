#Backstage Gateway Example

## Starting the Gateway
```golang
go run gateway_sample.go
```

## Unauthorized Request
```bash
curl -i http://tres.backstage.dev/xml
HTTP/1.1 401 Unauthorized
date: Thu, 14 May 2015 11:25:09 GMT
content-length: 22
content-type: text/plain; charset=utf-8
Connection: keep-alive

You must be logged in.
```

## Authorized Request
```bash
 curl -i http://tres.backstage.dev/xml -H "Authorization: secret"
HTTP/1.1 200 OK
alternate-protocol: 80:quic,p=1
content-length: 45
content-type: text/plain
date: Thu, 14 May 2015 11:25:50 GMT
server: Google Frontend
Connection: keep-alive

Hello, world! You have called me 2 times.
Foo
```
