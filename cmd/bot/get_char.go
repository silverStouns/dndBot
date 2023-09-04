package bot

import (
	"dndBot/internal/config"
	"dndBot/internal/database"
	"dndBot/internal/pkg/logger"
	"dndBot/internal/pkg/request"
	"dndBot/internal/pkg/session"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"path/filepath"
)

var GetCharHandler = func(update tgbotapi.Update) []tgbotapi.Chattable {
	//prepareChar()
	var arrayMessagesTextex []tgbotapi.Chattable
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

	req = request.NewRequest(getCharName, config.GetChar, update)

	req.SetMeta(request.MetaMsgText, update.Message.Text)
	switch req.Step {
	case getCharName:
		var resultText, imageUrl string
		var interator int
		charniki, err := database.GetCharInfo(DBconn, req.GetMeta(request.MetaTgID))
		if err != nil {
			logger.Error("GetCharInfo :%v", err)
		}
		for _, ch := range charniki {
			interator++
			imageUrl = ch.ImageUrl
			text := fmt.Sprintf(`Имя пользователя: %v
Имя модуля: %v
Имя персоонажа:%v
Расса:%v
Характеристики:%v
Опыт:%v
Уровень:%v
Класс:%v
Оружие:%+v
Навыки:%v
Бонус мастерства:%v
дополнительно усиление навыка:%+v
Золото:%v
Инвертарь:%v
Заклинания:%+v
Уникальные способности:%v
Уникальный ресурс:%v
Описание:%v
/////////////////////////////////
`, ch.NameUser, ch.NameModule, ch.NameChar, ch.Race, ch.Characteristic,
				ch.Experience, ch.Lvl, ch.Class, ch.Weapon, ch.Skills, ch.BonusMaster,
				ch.UnicueBonusSkills, ch.Gold, ch.Invertar, ch.Spels, ch.UnicueSpels, ch.Resurses,
				ch.Description)

			resultText += text
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, resultText)
		msg.ParseMode = "markdown"
		// Если перс 1, то выдать фотку
		if interator == 1 && len(imageUrl) > 0 {
			path := fmt.Sprintf("image/%v", imageUrl)
			absPath, err := filepath.Abs(path)
			if err != nil {
				logger.Error("abs", err)
			}

			// create the output file
			msgImage := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FilePath(absPath))
			arrayMessagesTextex = append(arrayMessagesTextex, msg)
			arrayMessagesTextex = append(arrayMessagesTextex, msgImage)
			return arrayMessagesTextex
		}

		arrayMessagesTextex = append(arrayMessagesTextex, msg)
		return arrayMessagesTextex
	}
	return []tgbotapi.Chattable{}
}
