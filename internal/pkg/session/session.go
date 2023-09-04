package session

import (
	"dndBot/internal/pkg/request"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var requestMap = make(map[int]request.Request)

var FinishSessionFlag = fmt.Errorf("finish session")

func GetSession(update tgbotapi.Update) (*request.Request, bool) {
	req, ok := requestMap[int(update.Message.From.ID)]
	return &req, ok
}

func SaveSession(update tgbotapi.Update, request *request.Request) {
	requestMap[int(update.Message.From.ID)] = *request
}

func FinishSession(update tgbotapi.Update) {
	delete(requestMap, int(update.Message.From.ID))
}
