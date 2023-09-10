package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/vladpi/film-exposure-bot/internal/bot"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	botToken := os.Getenv("BOT_TOKEN")

	bot, err := bot.NewBot(botToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Start(ctx)
}
