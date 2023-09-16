package bot

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
)

func NewBot(token string) (*bot.Bot, error) {
	opts := []bot.Option{
		bot.WithDebug(),
		bot.WithDefaultHandler(defautlHandler),
	}

	bot, err := bot.New(token, opts...)
	if err != nil {
		return nil, err
	}

	registerHandlers(bot)

	return bot, nil
}

func registerHandlers(b *bot.Bot) {
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, startHandler)
	b.RegisterHandlerMatchFunc(
		func(update *models.Update) bool {
			return len(update.Message.Photo) != 0
		},
		photoHandler,
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

func photoHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	filmKb := inline.New(b, inline.NoDeleteAfterClick())
	for _, f := range Films {
		filmKb = filmKb.Row().Button(f, []byte(f), onPhotoFilmSelect)
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

func onPhotoFilmSelect(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	film := string(data)

	shutterSpeedKb := inline.New(b, inline.NoDeleteAfterClick())
	for _, s := range ShutterSpeeds {
		shutterSpeedKb = shutterSpeedKb.Row().Button(s, []byte(s), onShutterSpeedSelect)
	}

	b.EditMessageCaption(ctx, &bot.EditMessageCaptionParams{
		ChatID:      mes.Chat.ID,
		MessageID:   mes.ID,
		Caption:     fmt.Sprintf("%s\n%s", mes.Caption, film),
		ReplyMarkup: shutterSpeedKb,
	})
}

func onShutterSpeedSelect(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	shutterSpeed := string(data)

	aperturesKb := inline.New(b, inline.NoDeleteAfterClick())
	for _, a := range Apertures {
		aperturesKb = aperturesKb.Row().Button(a, []byte(a), onApertureSelect)
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
