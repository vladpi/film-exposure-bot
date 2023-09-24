package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/vladpi/film-exposure-bot/internal/bot"
	"github.com/vladpi/film-exposure-bot/internal/repository"
	"github.com/vladpi/film-exposure-bot/internal/service"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	dbDSN := os.Getenv("DB_DSN")
	db, err := sqlx.Connect("sqlite3", dbDSN)
	if err != nil {
		log.Fatal(err)
	}

	filmRepo := repository.NewSQLFilmRepository(db)
	filmService := service.NewFilmService(filmRepo)

	botToken := os.Getenv("BOT_TOKEN")

	bot, err := bot.NewBot(botToken, filmService)
	if err != nil {
		log.Fatal(err)
	}

	bot.Start(ctx)
}
