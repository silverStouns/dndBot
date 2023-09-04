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

var AudioHandler = func(update tgbotapi.Update) []tgbotapi.Chattable {
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
		path := "image/what.mp3"
		rnd := rand.Intn(3)
		switch rnd {
		case 0:
			path = "image/what.mp3"
		case 1:
			path = "image/what2.mp3"
		case 2:
			path = "image/what3.mp3"

		}
		var arrayMessagesTextex []tgbotapi.Chattable
		msg := fmt.Sprintf("Звук негодования.")

		video := tgbotapi.NewAudio(req.GetChatID(), tgbotapi.FilePath(path))
		txt := tgbotapi.NewMessage(update.Message.Chat.ID, msg)
		txt.ParseMode = "markdown"
		arrayMessagesTextex = append(arrayMessagesTextex, txt)
		arrayMessagesTextex = append(arrayMessagesTextex, video)
		return arrayMessagesTextex
	}
	return []tgbotapi.Chattable{}
}
