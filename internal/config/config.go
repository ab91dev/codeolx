package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                string
	Env                 string
	DatabaseURL         string
	JWTKey              string
	StorageAccountID    string
	StorageAccessKey    string
	StorageAccessSecret string
	StorageBucket       string
}

func MustLoad() Config {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT is required")
	}

	env := os.Getenv("ENV")
	if env == "" {
		panic("ENV is required")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		panic("DATABASE_URL is required")
	}

	jwtKey := os.Getenv("JWT_KEY")
	if jwtKey == "" {
		panic("JWT_KEY is required")
	}

	storageAccountID := os.Getenv("STORAGE_ACCOUNT_ID")
	if storageAccountID == "" {
		panic("STORAGE_ACCOUNT_ID is required")
	}

	storageAccessKey := os.Getenv("STORAGE_ACCESS_KEY")
	if storageAccessKey == "" {
		panic("STORAGE_ACCESS_KEY is required")
	}

	storageAccessSecret := os.Getenv("STORAGE_ACCESS_SECRET")
	if storageAccessSecret == "" {
		panic("STORAGE_ACCESS_SECRET is required")
	}

	storageBucket := os.Getenv("BUCKET")
	if storageBucket == "" {
		panic("BUCKET is required")
	}

	return Config{
		Port:                port,
		Env:                 env,
		DatabaseURL:         dbURL,
		JWTKey:              jwtKey,
		StorageAccountID:    storageAccountID,
		StorageAccessKey:    storageAccessKey,
		StorageAccessSecret: storageAccessSecret,
		StorageBucket:       storageBucket,
	}
}
