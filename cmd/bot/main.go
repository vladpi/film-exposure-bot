package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/vladpi/film-exposure-bot/config"
	"github.com/vladpi/film-exposure-bot/internal/bot"
	"github.com/vladpi/film-exposure-bot/internal/repository"
	"github.com/vladpi/film-exposure-bot/internal/service"
)

func main() {
	config, err := config.LoadFromDotenv()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	db, err := sqlx.Connect("sqlite3", config.DBDSN)
	if err != nil {
		log.Fatal(err)
	}

	filmRepo := repository.NewSQLFilmRepository(db)
	filmService := service.NewFilmService(filmRepo)

	bot, err := bot.NewBot(config.BotToken, filmService)
	if err != nil {
		log.Fatal(err)
	}

	bot.Start(ctx)
}
