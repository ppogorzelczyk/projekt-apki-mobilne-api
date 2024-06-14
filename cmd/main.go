package main

import (
	"buymeagiftapi/internal/config"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	env := os.Getenv("ENV_TYPE")
	if env == "" {
		env = "development"
	}
	handlerOpt := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	logger := slog.New(slog.NewJSONHandler(os.Stderr, handlerOpt))
	slog.SetDefault(logger)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	slog.Info(fmt.Sprintf("Starting BuyMeAGift in %s mode", env))

	r := config.InitGin()

	r.Run()
}
