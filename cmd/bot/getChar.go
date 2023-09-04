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

// Тоже заглушка но чарника в будущем заполняется пользователем
//func prepareChar() {
//	characterisk := []Characteristic{
//		{
//			Name: "Сила",
//			Col:  20,
//			Mod:  0,
//		},
//		{
//			Name: "Ловкость",
//			Col:  14,
//			Mod:  0,
//		},
//		{
//			Name: "Телосложение",
//			Col:  16,
//			Mod:  0,
//		},
//		{
//			Name: "Интеллект",
//			Col:  12,
//			Mod:  0,
//		},
//		{
//			Name: "Мудрость",
//			Col:  10,
//			Mod:  0,
//		},
//		{
//			Name: "Харизма",
//			Col:  18,
//			Mod:  6,
//		},
//	}
//	weapon := []WeaponT{{
//		Upgrade:     1,
//		Description: "Идеально сделанный двуручный топор",
//		Type:        "Двуручный топор",
//		Damage:      "1d12",
//		UnicBonuses: "Очень тяжелый, минимально требуемая сила 16",
//	}, {
//		Upgrade:     0,
//		Description: "Длинный лук",
//		Type:        "Длинный лук",
//		Damage:      "1d8",
//	},
//	}
//	skils := []string{"Атлетика", "Медицина"}
//	characterisk = ResultModifiCharacteristic(characterisk)
//	Charnins = make(map[string]Char)
//	unicSpel := map[string]string{"Подсечка": "Сбивает с ног сложность спаса 12"}
//	unicSkills := map[string]string{"Атлетика": "10"}
//	f := Char{
//		NameUser:          "Anton",
//		NameModule:        "TestModule",
//		NameChar:          "Абоба",
//		Race:              "Человек",
//		Experience:        900,
//		Lvl:               ResultExpLvl(900),
//		Class:             "Воин/Мастер боя",
//		Characteristic:    characterisk,
//		Weapon:            weapon,
//		Skills:            skils,
//		BonusMaster:       ResultBonusMaster(3),
//		UnicueBonusSkills: unicSkills,
//		Gold:              10,
//		Invertar: `Тут лежат мои пожитки
//Еда на 10 приёмов
//Вода на 5 приёмов
//20 факелов
//100 футов верёвки`,
//		Spels:       []Spels{},
//		UnicueSpels: unicSpel,
//		Resurses:    "Кости превосходства d8, 4шт",
//		Description: "Обычный воин",
//		ImageUrl:    fmt.Sprintf("image/%v.png", 22),
//	}
//	Charnins["Абоба"] = f
//}
