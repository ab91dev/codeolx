package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/ab91dev/codeolx/internal/config"
	"github.com/ab91dev/codeolx/internal/db"
	"github.com/ab91dev/codeolx/internal/handlers"
	"github.com/ab91dev/codeolx/internal/middleware"
	"github.com/ab91dev/codeolx/internal/storage"
)

func main() {
	cfg := config.MustLoad()

	db, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("main.db.connect: %v", err)
	}

	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	fmt.Println("database connected")

	// storage initialisation
	store, err := storage.NewR2Storage(context.TODO(), storage.R2Config{
		AccountID:    cfg.StorageAccountID,
		AccessKey:    cfg.StorageAccessKey,
		AccessSecret: cfg.StorageAccessSecret,
		Bucket:       cfg.StorageBucket,
	})
	if err != nil {
		log.Fatalf("main.storage.r2", err)
	}

	fmt.Println("storage initialised...")
	fmt.Println("starting the server...")

	listingsHandler := handlers.NewListingHandlerParams(db, logger, store)
	authHandler := handlers.NewAuthHandler(db, logger, cfg)
	requireAuth := middleware.RequireAuth(logger, cfg.JWTKey)
	uploadHandler := handlers.NewUploadHandler(logger, store)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handlers.Healthz)
	mux.HandleFunc("GET /listings", listingsHandler.List)
	mux.Handle("DELETE /listings/{id}", requireAuth(http.HandlerFunc(listingsHandler.Delete)))
	mux.Handle("POST /listings", requireAuth(http.HandlerFunc(listingsHandler.Create)))
	mux.HandleFunc("POST /signup", authHandler.SignUp)
	mux.HandleFunc("POST /signin", authHandler.SignIn)
	mux.Handle("POST /uploads/presign", requireAuth(http.HandlerFunc(uploadHandler.Presign)))

	handler := middleware.RequestId(mux)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 40,
		IdleTimeout:  time.Second * 120,
	}

	log.Printf("Server is listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
