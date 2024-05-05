package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"jeisaraja/websocket_learn/database"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

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
	flag.StringVar(&cfg.db.dsn, "dh-dsn", "postgres://chat_ws:password@localhost/chat_ws", "DB connection string")
	flag.Parse()
	setupAPI(&cfg)
	log.Println("Listening on port 5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}

func setupAPI(cfg *config) {

	var Queries = database.Queries{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db, err := connectDB(cfg, ctx)

	if err != nil {
		log.Fatal(err.Error())
	}

	Queries.DB = db
	manager := NewManager(ctx, &Queries)
	defer db.Close()
	fmt.Printf("database connection pool established")
	http.HandleFunc("/", handleNotFound)
	http.HandleFunc("/ws", manager.serveWs)
	http.HandleFunc("/login", manager.login)
	http.HandleFunc("/signup", manager.signup)
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

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	msg := "Nothing found"
	w.Write([]byte(msg))
}
