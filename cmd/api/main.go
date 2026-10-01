package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ab91dev/codeolx/internal/config"
	"github.com/ab91dev/codeolx/internal/db"
	"github.com/ab91dev/codeolx/internal/handlers"
)

func main() {
	cfg := config.MustLoad()

	db, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("main.db.connect: %v",err)
	}

	fmt.Println("database connected")
	fmt.Println("starting the server on render")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handlers.Healthz)
	mux.HandleFunc("GET /listings", handlers.Listings(db))

	srv := &http.Server{
		Addr: ":" + cfg.Port,
		Handler: mux,
		ReadTimeout: time.Second * 10,
		WriteTimeout: time.Second * 40,
		IdleTimeout: time.Second * 120,
	}

	log.Printf("Server is listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil{
		log.Fatalf("Server failed: %v",err)
	}
}
