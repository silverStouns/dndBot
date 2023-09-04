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

var ModulesInformation map[string]ModuleInfo

type ModuleInfo struct {
	DungeonMaster string // Имя ведущего
	Num           int    // Номер модуля
	//Gamers        map[string]string // Имя игрока и имя персоонажа
	//Col           int               // Количество игроков(суммарно играющих в модуль)
	Description string //Описание модуля(должно заполнятся мастером, максимум 32тыс символов(пожалуйста не хуярьте сюда войну и мир оч прошу))
}

var ModulInfoHandler = func(update tgbotapi.Update) []tgbotapi.Chattable {
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

	req = request.NewRequest(modulInfoStep, config.ModulInfo, update)

	req.SetMeta(request.MetaMsgText, update.Message.Text)
	switch req.Step {
	case modulInfoStep:
		testPrepareInfo()
		var msg string
		for name, infoModule := range ModulesInformation {
			//var gamers string
			msg += fmt.Sprintf(`Название: %v
Номер: %v
Кто ведёт: %v
Описание: %v
`, name, infoModule.Num, infoModule.DungeonMaster, infoModule.Description)
			//for name, class := range infoModule.Gamers {
			//	gamers += fmt.Sprintf("%v/%v\n", name, class)
			//}
			//msg += fmt.Sprintf("\n\n*%v* \nКто ведёт:%v\nИгроки:\n%vКолличество:%v\nОписание:%v", name, infoModule.DungeonMaster, gamers, infoModule.Col, infoModule.Description)
		}

		msgTg := tgbotapi.NewMessage(req.GetChatID(), msg)
		msgTg.ParseMode = "markdown"
		msgTg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true, Selective: true}
		return []tgbotapi.Chattable{msgTg}
	}
	return []tgbotapi.Chattable{}
}

// в будущем такая штука должна заполнятся сама это заглушка
func testPrepareInfo() {
	ModulesInformation = make(map[string]ModuleInfo)
	ModulesInformation["Пещеры Мутаций"] = ModuleInfo{
		DungeonMaster: "Антон",
		Num:           1,
		Description:   "Пещеры где вы подвергаетесь постоянным изменениями, какие же тайны они скрывают?",
	}

	ModulesInformation["Безбожный мир"] = ModuleInfo{
		DungeonMaster: "Димас",
		Num:           2,
		Description:   "Хуй пойми чё тут происходит(димас заполни плиз пидорас)))",
	}
}
