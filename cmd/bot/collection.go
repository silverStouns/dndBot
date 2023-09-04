package bot

import (
	"dndBot/internal/config"
	"dndBot/internal/pkg/logger"
	"dndBot/internal/pkg/request"
	"dndBot/internal/pkg/session"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var infoGame = []string{"Буду", "Не буду", "Бан!"}
var CollectionHandler = func(update tgbotapi.Update) []tgbotapi.Chattable {
	var err error
	req, _ := session.GetSession(update)
	defer func() {
		if err == nil {
			session.SaveSession(update, req)
		}
		if errors.Is(err, session.FinishSessionFlag) {
			logger.Error("err:%v", err)
			session.FinishSession(update)
		}
		if err != nil {
			logger.Error("err: %v", err)
		}
	}()

	req = request.NewRequest(collectionStep, config.Collection, update)

	req.SetMeta(request.MetaMsgText, update.Message.Text)
	switch req.Step {
	case collectionStep:
		poll := tgbotapi.NewPoll(req.GetChatID(), "Кто будет на игре?", infoGame...)
		poll.IsAnonymous = false

		return []tgbotapi.Chattable{poll}
	}
	return []tgbotapi.Chattable{}
}
