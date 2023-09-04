package bot

import (
	"dndBot/internal/config"
	"dndBot/internal/database"
	"dndBot/internal/pkg/logger"
	"dndBot/internal/pkg/request"
	"dndBot/internal/pkg/session"
	"encoding/json"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Id шагов по созданию чарника
const (
	GetCharPrimer = iota + 99 // Выдаёт пример который пользователь должен заполнить
	SetCharInfo
)

// В будущем я научу тебя создавать чарник
var CreateCharHandler = func(update tgbotapi.Update) []tgbotapi.Chattable {
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
		req = request.NewRequest(GetCharPrimer, config.Createchar, update)
	}
	req.SetMeta(request.MetaMsgText, update.Message.Text)
	switch req.Step {
	case GetCharPrimer:
		tmpl := getTemplate()
		//t, err := json.Marshal(tmpl)
		//if err != nil {
		//	logger.Error("err:%v", err)
		//}
		msgTg := tgbotapi.NewMessage(req.GetChatID(), tmpl)
		//msgTg.ParseMode = "markdown"
		msgTg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true, Selective: true}
		req.Step = SetCharInfo
		return []tgbotapi.Chattable{msgTg}
	case SetCharInfo:
		ch := database.Char{}
		txt := "Произошла ошибка, обратитесь к разработчику"
		err := json.Unmarshal([]byte(update.Message.Text), &ch)
		if err != nil {
			logger.Error("err:%v", err)
		}
		userID, err := database.GetGamerId(DBconn, req.GetMeta(request.MetaTgID))
		if err != nil {
			logger.Error("err:%v", err)
		}
		err = database.CreateChar(DBconn, ch, userID)
		if err != nil {
			logger.Error("err:%v", err)
		}

		if err == nil {
			txt = "Персонаж добавлен"
		}

		msgTg := tgbotapi.NewMessage(req.GetChatID(), txt)
		//msgTg.ParseMode = "markdown"
		msgTg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true, Selective: true}
		req.Step = SetCharInfo
		return []tgbotapi.Chattable{msgTg}
	}
	return []tgbotapi.Chattable{}
}

func getTemplate() string {
	byteeee, _ := json.MarshalIndent(map[string]interface{}{
		"name_char": "Имя персонажа",
		"characteristic": []map[string]interface{}{
			{"name": "Сила", "col": 10},
			{"name": "Ловкость", "col": 10},
			{"name": "Телосложение", "col": 10},
			{"name": "Мудрость", "col": 10},
			{"name": "Интеллект", "col": 10},
			{"name": "Харизма", "col": 10},
		},
		"race":       "Раса",
		"experience": 900,
	}, "\t", "")
	return string(byteeee)
	//"experience":900,
	//"class":"Класс персонажа/подкласс(если есть)",
	//"weapon":[{"upgrade":1,"description":"Описание оружия","type":"Тип оружия","damage":"Урон оружия например 1d10","unic_bonuses":"Уникальный бонус оружия, например возможность кровотечения(желательно указать сложность)"},
	//{"upgrade":2,"description":"Описание оружия2","type":"Тип оружия2","damage":"Урон оружия например 1d10","unic_bonuses":"Уникальный бонус оружия, например возможность кровотечения(желательно указать сложность)"}],
	//"skills":["Атлетика","Акробатика","Магия"],
	//"unicue_bonus_skills":{"Акробатика":"10","Атлетика":"10"},
	//"gold":500,
	//"invertar":"Опишите что лежит в вашем инвертаре(в это поле вмещается 32к символов)",
	//"spels":[{"name":"Огненный снаряд","lvl":0,"damage":"1d10","type_spas":"Нет(если есть укажите какой тип)","hard_spas":0}],
	//"unicue_spels":{"Подсечка":"Сбивает с ног, сложность 14"},
	//"resurses":"Уникальный рессурс, например: Ярость 2 до отдыха, Возложение рук 20",
	//"description":"Описание вашего персоонажа",
	//"image_url":"Если хотим картинку, то кидаем её в личку антону и говорим имя перса",
	//"num_module":1}`
}

//func getTemplate() database.Char {
//	tmpl := database.Char{
//		NameUser:   "Имя пользователя(заполнять не требуется)",
//		NameModule: "Имя модуля(заполнять не требуется)",
//		NameChar:   "Имя персонажа",
//		Race:       "Раса",
//		Characteristic: []database.Characteristic{
//			{
//				Name: "Сила",
//				Col:  10,
//			},
//			{
//				Name: "Ловкость",
//				Col:  10,
//			},
//			{
//				Name: "Телосложение",
//				Col:  10,
//			},
//			{
//				Name: "Мудрость",
//				Col:  10,
//			},
//			{
//				Name: "Интеллект",
//				Col:  10,
//			},
//			{
//				Name: "Харизма",
//				Col:  10,
//			},
//		},
//		Experience: 900,
//		Class:      "Класс персонажа/подкласс(если есть)",
//		Weapon: []database.WeaponT{
//			{
//				Upgrade:     1,
//				Description: "Описание оружия",
//				Type:        "Тип оружия",
//				Damage:      "Урон оружия например 1d10",
//				UnicBonuses: "Уникальный бонус оружия, например возможность кровотечения(желательно указать сложность)",
//			},
//		},
//		Skills:            []string{"Атлетика", "Акробатика", "Магия"},
//		UnicueBonusSkills: map[string]string{"Акробатика": "10"},
//		Gold:              500,
//		Invertar:          "Опишите что лежит в вашем инвертаре(в это поле вмещается 32к символов)",
//		Spels: []database.Spels{
//			{
//				Name:     "Огненный снаряд",
//				Lvl:      0,
//				Damage:   "1d10",
//				TypeSpas: "Нет(если есть укажите какой тип)",
//				HardSpas: 0,
//			},
//		},
//		UnicueSpels: map[string]string{"Подсечка": "Сбивает с ног, сложность 14"},
//		Resurses:    "Уникальный рессурс, например: Ярость 2 до отдыха, Возложение рук 20",
//		Description: "Описание вашего персоонажа",
//		ImageUrl:    "Если хотим картинку, то кидаем её в личку антону и говорим имя перса",
//		NumModule:   1,
//	}
//	return tmpl
//}
