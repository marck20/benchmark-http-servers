FROM golang

WORKDIR /goab

COPY goab_mutex.go .

RUN go get github.com/valyala/fasthttp

RUN go build -o goab_mutex

ENTRYPOINT ["/goab/goab_mutex"]


