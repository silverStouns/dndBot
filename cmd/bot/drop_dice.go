package bot

import (
	"dndBot/internal/pkg/logger"
	"dndBot/internal/pkg/request"
	"dndBot/internal/pkg/session"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math/rand"
)

var DropDiceHandler = func(update tgbotapi.Update) []tgbotapi.Chattable {
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
	req = request.NewRequest(1, "/drop_dice", update)

	req.SetMeta(request.MetaMsgText, update.Message.Text)
	switch req.Step {
	case 1:
		msg := ""
		drop := rand.Intn(21)
		switch drop {
		case 0:
			var arrayMessagesTextex []tgbotapi.Chattable
			msg = fmt.Sprintf("Ты бросил d20 на: 1, ну и лох...")

			video := tgbotapi.NewVideo(req.GetChatID(), tgbotapi.FilePath("image/2yhkuz.mp4"))
			txt := tgbotapi.NewMessage(update.Message.Chat.ID, msg)
			arrayMessagesTextex = append(arrayMessagesTextex, txt)
			arrayMessagesTextex = append(arrayMessagesTextex, video)
			return arrayMessagesTextex
		case 1:
			msg = fmt.Sprintf("Ты бросил d20 на: %v, ну и лох...", drop)
		case 20:
			msg = fmt.Sprintf("Ты бросил d20 на: %v, заебумба", drop)
		case 21:
			msg = fmt.Sprintf("Ты бросил d20 на: 20, заебумба")
		default:
			msg = fmt.Sprintf("Ты бросил d20 на: %v", drop)
		}

		msgTg := tgbotapi.NewMessage(req.GetChatID(), msg)
		msgTg.ParseMode = "markdown"
		msgTg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true, Selective: true}
		return []tgbotapi.Chattable{msgTg}
	}
	return []tgbotapi.Chattable{}
}
