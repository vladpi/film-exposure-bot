package bot

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
	"github.com/vladpi/film-exposure-bot/internal/domain/film"
)

func NewBot(token string, fs film.Service) (*bot.Bot, error) {
	opts := []bot.Option{
		bot.WithDebug(),
		bot.WithDefaultHandler(defautlHandler),
	}

	bot, err := bot.New(token, opts...)
	if err != nil {
		return nil, err
	}

	registerHandlers(bot, fs)

	return bot, nil
}

func registerHandlers(b *bot.Bot, fs film.Service) {
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, startHandler)
	b.RegisterHandlerMatchFunc(
		func(update *models.Update) bool {
			return len(update.Message.Photo) != 0
		},
		photoHandler(fs),
	)
}

func defautlHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    update.Message.Chat.ID,
		MessageID: update.Message.ID,
	})
}

func startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Привет!\nЯ бот, который умеет записывать параметры съемки при фотографировании.\n\nОтправь мне фотографию, и посмотри как это работает 😉",
	})
}

func photoHandler(fs film.Service) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		films, err := fs.GetAll()
		if err != nil {
			log.Fatal(err) // FIXME добавить правильную обработку ошибок
		}

		filmKb := inline.New(b, inline.NoDeleteAfterClick())
		for _, f := range films {
			filmKb = filmKb.Row().Button(f.Name, []byte(strconv.FormatInt(f.ID, 10)), onPhotoFilmSelect(fs))
		}

		photo := update.Message.Photo[len(update.Message.Photo)-1]

		b.SendPhoto(ctx, &bot.SendPhotoParams{
			ChatID:      update.Message.Chat.ID,
			Photo:       &models.InputFileString{Data: photo.FileID},
			ReplyMarkup: filmKb,
		})
		b.DeleteMessage(ctx, &bot.DeleteMessageParams{
			ChatID:    update.Message.Chat.ID,
			MessageID: update.Message.ID,
		})
	}
}

func onPhotoFilmSelect(fs film.Service) inline.OnSelect {
	return func(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
		filmId, err := strconv.ParseInt(string(data), 10, 0)
		if err != nil {
			log.Fatal(err) // FIXME добавить правильную обработку ошибок
		}
		film, err := fs.Get(filmId)
		if err != nil {
			log.Fatal(err) // FIXME добавить правильную обработку ошибок
		}

		shutterSpeedKb := inline.New(b, inline.NoDeleteAfterClick())
		for i, s := range ShutterSpeeds {
			if i%3 == 0 {
				shutterSpeedKb = shutterSpeedKb.Row()
			}
			shutterSpeedKb = shutterSpeedKb.Button(s, []byte(s), onShutterSpeedSelect)
		}

		b.EditMessageCaption(ctx, &bot.EditMessageCaptionParams{
			ChatID:      mes.Chat.ID,
			MessageID:   mes.ID,
			Caption:     fmt.Sprintf("%s\n%s", mes.Caption, film.Name),
			ReplyMarkup: shutterSpeedKb,
		})
	}
}

func onShutterSpeedSelect(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	shutterSpeed := string(data)

	aperturesKb := inline.New(b, inline.NoDeleteAfterClick())
	for i, a := range Apertures {
		if i%3 == 0 {
			aperturesKb = aperturesKb.Row()
		}
		aperturesKb = aperturesKb.Button(a, []byte(a), onApertureSelect)
	}

	b.EditMessageCaption(ctx, &bot.EditMessageCaptionParams{
		ChatID:      mes.Chat.ID,
		MessageID:   mes.ID,
		Caption:     fmt.Sprintf("%s\n%s", mes.Caption, shutterSpeed),
		ReplyMarkup: aperturesKb,
	})
}

func onApertureSelect(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	aperture := string(data)

	b.EditMessageCaption(ctx, &bot.EditMessageCaptionParams{
		ChatID:      mes.Chat.ID,
		MessageID:   mes.ID,
		Caption:     fmt.Sprintf("%s 𝒇%s", mes.Caption, aperture),
		ReplyMarkup: &models.ReplyKeyboardRemove{},
	})
}
