package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ab91dev/codeolx/internal/config"
	"github.com/ab91dev/codeolx/internal/handlers"
)

func main() {
	cfg := config.MustLoad()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handlers.Healthz)

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
