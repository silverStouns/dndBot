package bot

import (
	"dndBot/internal/config"
	"dndBot/internal/pkg/logger"
	"dndBot/internal/pkg/request"
	"dndBot/internal/pkg/session"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	startStep = iota
	updateStep
	summonStep
	helpStep
	collectionStep
	modulInfoStep
	getCharName
	getCharInfo
)

var StartHandler = func(update tgbotapi.Update) []tgbotapi.Chattable {
	var err error
	req, continueHandling := session.GetSession(update)
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
	if !continueHandling {
		req = request.NewRequest(startStep, config.Start, update)
	}
	req.SetMeta(request.MetaMsgText, update.Message.Text)
	switch req.Step {
	case startStep:
		msg := tgbotapi.NewMessage(req.GetChatID(), "Приветствую авантюрист!")
		msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true, Selective: true}
		req.Step = updateStep
		return []tgbotapi.Chattable{msg}
	default:
		msg := tgbotapi.NewMessage(req.GetChatID(), "С возвращением авантюрист!")
		msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true, Selective: true}
		return []tgbotapi.Chattable{msg}
	}
}
