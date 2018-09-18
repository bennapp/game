FROM golang:1.10

WORKDIR /go/src/github.com/bennapp/game

RUN go get -u github.com/nsf/termbox-go && \
    go get github.com/go-redis/redis && \
    go get github.com/gorilla/websocket && \
    go get github.com/google/uuid && \
    go get github.com/vmihailenco/msgpack

COPY back-end back-end

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
    -o ./zerozero-linux-amd64 back-end/websocket_server/websocket_server.go

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
    -o ./zerozero-new-world-linux-amd64 back-end/new_world/new_world.go

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
    -o ./zerozero-spawn-coins-linux-amd64 back-end/spawn_coins/spawn_coins.go

FROM alpine:3.7

COPY --from=0 /go/src/github.com/bennapp/game/zerozero-linux-amd64 /usr/bin/zerozero
COPY --from=0 /go/src/github.com/bennapp/game/zerozero-new-world-linux-amd64 /usr/bin/zerozero-new-world
COPY --from=0 /go/src/github.com/bennapp/game/zerozero-spawn-coins-linux-amd64 /usr/bin/zerozero-spawn-coins
COPY game-config /go/src/github.com/bennapp/game/game-config

CMD ["/usr/bin/zerozero"]
