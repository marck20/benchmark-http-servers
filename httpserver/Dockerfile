FROM golang

WORKDIR /httpserver

COPY httpserver.go .

RUN go get github.com/valyala/fasthttp

RUN go build -o httpserver

ENTRYPOINT ["/httpserver/httpserver"]
