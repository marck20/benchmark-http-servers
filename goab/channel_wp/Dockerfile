FROM golang

WORKDIR /goab

COPY goab_workerpool_channel.go .

RUN go get github.com/valyala/fasthttp

RUN go build -o goab_workerpool_channel

ENTRYPOINT ["/goab/goab_workerpool_channel"]


