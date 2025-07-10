package main

import (
	"log"
	"net/http"

	"github.com/Udehlee/strim/internals"
)

func main() {
	hub := internals.NewHub()
	go hub.Run()

	wsConn := internals.NewWsConn(hub)
	mux := wsConn.Route()

	log.Println("Starting WebSocket signaling server started on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("failed to start server")

	}
}
