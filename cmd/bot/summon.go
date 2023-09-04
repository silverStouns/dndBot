package bot

import (
	"dndBot/internal/config"
	"dndBot/internal/pkg/logger"
	"dndBot/internal/pkg/request"
	"dndBot/internal/pkg/session"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var SummonHandler = func(update tgbotapi.Update) []tgbotapi.Chattable {
	//prepareActiveGamers()
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
	req = request.NewRequest(summonStep, config.Summon, update)

	req.SetMeta(request.MetaMsgText, update.Message.Text)
	switch req.Step {
	case summonStep:
		message := ""
		// формируем ссылку
		for id, name := range Users {
			nameString := fmt.Sprintf("[%v](tg://user?id=%v)", name, id)
			message += fmt.Sprintf("%v ", nameString)
		}
		// встраиваем

		msg := tgbotapi.NewMessage(req.GetChatID(), message)
		msg.ParseMode = "MarkdownV2"
		msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true, Selective: true}
		return []tgbotapi.Chattable{msg}
	}
	return []tgbotapi.Chattable{}
}

//func prepareActiveGamers() {
//	// Тут надо будет парсить файл с именами
//	Users = []string{"@immortalMage", "@Nereeon"} // це заглушка
//	for _, gamer := range Users {
//		activeGamers += fmt.Sprintf(" %v", gamer)
//	}
//}
