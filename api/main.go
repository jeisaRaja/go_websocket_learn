package main

import (
	"log"
	"net/http"
)

func main() {
	setupAPI()
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func setupAPI() {
	manager := NewManager()

	http.Handle("/", http.FileServer(http.Dir("../client")))
	http.HandleFunc("/ws", manager.serveWs)
}
