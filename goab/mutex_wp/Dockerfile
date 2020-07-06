FROM golang

WORKDIR /goab

COPY goab_workerpool_mutex.go .

RUN go get github.com/valyala/fasthttp

RUN go build -o goab_workerpool_mutex

ENTRYPOINT ["/goab/goab_workerpool_mutex"]


