package bot

import (
	"dndBot/internal/pkg/logger"
	"dndBot/internal/pkg/request"
	"dndBot/internal/pkg/session"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math/rand"
)

var VideoHandler = func(update tgbotapi.Update) []tgbotapi.Chattable {
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
	req = request.NewRequest(1, "/not_push", update)

	req.SetMeta(request.MetaMsgText, update.Message.Text)
	switch req.Step {
	case 1:
		videos := "image/2xp1kc.mp4"
		rnd := rand.Intn(4)
		switch rnd {
		case 0:
			videos = "image/30fku1.mp4"
		case 1:
			videos = "image/2ujbar.mp4"
		case 2:
			videos = "image/2dmixx.mp4"
		}
		video := tgbotapi.NewVideo(req.GetChatID(), tgbotapi.FilePath(videos))
		return []tgbotapi.Chattable{video}
	}
	return []tgbotapi.Chattable{}
}
