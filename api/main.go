package main

import (
	"context"
	"log"
	"net/http"
)

type UserAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	setupAPI()
	log.Fatal(http.ListenAndServeTLS(":3000", "server.crt", "server.key", nil))
}

func setupAPI() {

	ctx := context.Background()
	manager := NewManager(ctx)

	http.Handle("/", http.FileServer(http.Dir("../client")))
	http.HandleFunc("/ws", manager.serveWs)
	http.HandleFunc("/login", manager.login)
}
