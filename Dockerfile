FROM golang:1.10

WORKDIR /go/src/github.com/bennapp/game

RUN go get -u github.com/nsf/termbox-go && \
    go get github.com/go-redis/redis && \
    go get github.com/gorilla/websocket

COPY back-end back-end

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
    -o ./zerozero-linux-amd64 back-end/websocketServer.go

FROM alpine:3.7

COPY --from=0 /go/src/github.com/bennapp/game/zerozero-linux-amd64 /usr/bin/zerozero

CMD ["/usr/bin/zerozero"]
