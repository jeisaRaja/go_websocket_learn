package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

type UserAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development | staging | production)")
	flag.StringVar(&cfg.db.dsn, "dh-dsn", "postgres://go_websocket_learn:secret@localhost/go_websocket_learn", "DB connection string")
	flag.Parse()
	setupAPI(&cfg)
	log.Fatal(http.ListenAndServeTLS(":3000", "server.crt", "server.key", nil))
}

func setupAPI(cfg *config) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db, err := connectDB(cfg, ctx)

	if err != nil {
		log.Fatal(err.Error())
	}

	manager := NewManager(ctx, db)
	defer db.Close()
	fmt.Printf("database connection pool established")
	http.Handle("/", http.FileServer(http.Dir("../client")))
	http.HandleFunc("/ws", manager.serveWs)
	http.HandleFunc("/login", manager.login)
}

func connectDB(cfg *config, ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
