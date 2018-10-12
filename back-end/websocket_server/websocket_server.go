package main

import (
	"../ws"
	"log"
	"net/http"
)

import _ "net/http/pprof"

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	ws.RunServer()
}
