package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	Env string
	Goos string
	Goarch string
}

func MustLoad() Config {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == ""{
		panic("PORT is required")
	}

	env := os.Getenv("ENV")
	if env == ""{
		panic("ENV is required")
	}

	goos := os.Getenv("GOOS")
	if goos == ""{
		panic("GOOS is required")
	}

	goarch := os.Getenv("GOARCH")
	if goarch == ""{
		panic("GOARCH is required")
	}

	return Config{
		Port: port,
		Env: env,
		Goos: goos,
		Goarch: goarch,
	}
}
