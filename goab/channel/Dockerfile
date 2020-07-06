FROM golang

WORKDIR /goab

COPY goab_channel.go .

RUN go get github.com/valyala/fasthttp

RUN go build -o goab_channel

ENTRYPOINT ["/goab/goab_channel"]


