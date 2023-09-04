package user_token

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func PrepareUserToken(update tgbotapi.Update) string {
	return fmt.Sprintf("%d:%d", update.Message.From.ID, update.Message.Chat.ID)
}
