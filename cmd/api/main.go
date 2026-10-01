package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
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

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level: slog.LevelDebug,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	fmt.Println("database connected")
	fmt.Println("starting the server...")

	listingsHandler := handlers.NewListingHandlerParams(db,logger)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handlers.Healthz)
	mux.HandleFunc("GET /listings", listingsHandler.List)
	mux.HandleFunc("DELETE /listings/{id}", listingsHandler.Delete)
	
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
