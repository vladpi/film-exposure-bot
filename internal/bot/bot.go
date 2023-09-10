package bot

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
)

func NewBot(token string) (*bot.Bot, error) {
	opts := []bot.Option{
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
}

func defautlHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    update.Message.Chat.ID,
		MessageID: update.Message.ID,
	})
}

func startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	kb := inline.New(b).Row().Button("Cancel", []byte("cancel"), onInlineKeyboardSelect)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Hello, world!",
		ReplyMarkup: kb,
	})
}

func onInlineKeyboardSelect(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Chat.ID,
		Text:   "You selected: " + string(data),
	})
}
