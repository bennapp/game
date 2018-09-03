package ws

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8081", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)

	//if r.URL.Path == "/bundle.js" {
	//	http.ServeFile(w, r, "../front-end/dist/bundle.js")
	//	return
	//}
	//
	//if r.URL.Path == "/index.js" {
	//	http.ServeFile(w, r, "../front-end/src/client/index.js")
	//	return
	//}

	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	//http.ServeFile(w, r, "../front-end/dist/index.html")
	//http.ServeFile(w, r, "../front-end/public/index.html")
}

func RunServer() {
	flag.Parse()
	hub := newHub()
	go hub.run()

	// Clean this up after webpack assets refactor
	//fs := http.FileServer(http.Dir("../front-end/assets"))
	//http.Handle("/assets/", http.StripPrefix("/assets", fs))

	//http.HandleFunc("/", serveHome)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
