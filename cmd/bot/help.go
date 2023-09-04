package bot

import (
	"dndBot/internal/config"
	"dndBot/internal/pkg/logger"
	"dndBot/internal/pkg/request"
	"dndBot/internal/pkg/session"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var HelpHandler = func(update tgbotapi.Update) []tgbotapi.Chattable {
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

	req = request.NewRequest(helpStep, config.Help, update)

	req.SetMeta(request.MetaMsgText, update.Message.Text)
	switch req.Step {
	case helpStep:
		msg := tgbotapi.NewMessage(req.GetChatID(), `Вот мои команды:
/help - Предоставить информацию о командах
/Summon - Призвать всех активных игроков
/Collection - Вызвать опрос о том кто будет на игре
/moduls - Информация о идущих модулях
/get_char - Информация о твоём персоонаже
/drop_dice - бросает кубик d20
/what - нажать в случае единички на кубе!
/not_push - не советую нажимать -_о`)

		msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true, Selective: true}
		return []tgbotapi.Chattable{msg}
	}
	return []tgbotapi.Chattable{}
}
