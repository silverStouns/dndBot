package database

import (
	"database/sql"
	"dndBot/internal/pkg/logger"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

// Вытягивает всю инфу по всем персам игрока(Имена сука уникальными надо сделать!!!)
func GetCharInfo(conn *DBConnector, tgid string) (chrkini map[string]Char, err error) {
	chrkini = make(map[string]Char)
	logger.Debug("tg id:%v", tgid)
	intTgId, err := strconv.Atoi(tgid)
	if err != nil {
		logger.Error("err:%v", err)
	}
	var (
		userName, moduleName, class, weapon, charsName, characteristic, inveroty, spels, skills sql.NullString
		exp, gold                                                                               sql.NullInt64
		unicSkills, unicSpels, unicResurses, imageUrl, description, race                        sql.NullString
	)
	rows, err := conn.Connector.Query(`select u.name, m.module_name,
       main.chars.name,race, characteristic, exp,
       class, weapon, skills, unic_skills,
       gold, invertory, spels,
       unic_spels, unic_resurses,
       image_url, chars.desription
from chars
inner join main.moduls m on m.id = chars.module_id
inner join main.user u on u.id = chars.gamer_id
where u.telegram_id = $1`, intTgId)
	if err != nil {
		logger.Error("err:%v", err)
	}
	for rows.Next() {
		err = rows.Scan(&userName, &moduleName, &charsName, &race, &characteristic, &exp, &class, &weapon, &skills, &unicSkills, &gold, &inveroty, &spels,
			&unicSpels, &unicResurses, &imageUrl, &description)
		if err != nil {
			logger.Error("err:%v", err)
		}
		ch := Char{
			NameUser:          userName.String,
			NameModule:        moduleName.String,
			NameChar:          charsName.String,
			Race:              race.String,
			Characteristic:    ResultModifiCharacteristic(convertCharacteriscs(characteristic.String)),
			Experience:        exp.Int64,
			Lvl:               ResultExpLvl(exp.Int64),
			Class:             class.String,
			Weapon:            convertWeapon(weapon.String),
			Skills:            converSkils(skills.String),
			BonusMaster:       ResultBonusMaster(ResultExpLvl(exp.Int64)),
			UnicueBonusSkills: converUnicBonusSkills(unicSkills.String),
			Gold:              gold.Int64,
			Invertar:          inveroty.String,
			Spels:             convertSpels(spels.String),
			UnicueSpels:       converUnicSples(unicSpels.String),
			Resurses:          unicResurses.String,
			Description:       description.String,
			ImageUrl:          imageUrl.String,
		}
		chrkini[ch.NameChar] = ch
	}

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, nil
	case err != nil:
		logger.Debug("GetModulsInfo err")
		return nil, err
	default:
		return chrkini, nil
	}
}

// Функции конверторы(вынести в отдельный пакет(возможно сделать методами char))
func convertCharacteriscs(data string) (res []Characteristic) {
	if err := json.Unmarshal([]byte(data), &res); err != nil {
		logger.Error("err:%v", err)
	}
	return
}
func convertSpels(data string) (res []Spels) {
	if err := json.Unmarshal([]byte(data), &res); err != nil {
		logger.Error("err:%v", err)
	}
	return
}
func convertWeapon(data string) (res []WeaponT) {
	if err := json.Unmarshal([]byte(data), &res); err != nil {
		logger.Error("err:%v", err)
	}
	return
}

// Конвертируем строку в массив навыков(нужно что бы в будущем их считать!)
func converSkils(skills string) []string {
	result := strings.Split(skills, ", ")
	return result
}

// Конвертируем строку в map где ключ название а значение эффект(подумать нужно ли так делать? может просто строка...)
func converUnicSples(skills string) map[string]string {
	result := make(map[string]string)
	if len(skills) == 0 {
		return result
	}
	// В данном массиве значение уже в виде Key:value
	skils := strings.Split(skills, ":::")
	for _, sk := range skils {
		f := strings.Split(sk, ":")
		result[f[0]] = f[1]
	}
	return result
}

// Конвертируем строку в map где ключ название а значение эффект(подумать нужно ли так делать? может просто строка...)
func converUnicBonusSkills(spels string) map[string]string {
	result := make(map[string]string)
	if len(spels) == 0 {
		return result
	}
	// В данном массиве значение уже в виде Key:value
	skils := strings.Split(spels, ",")
	for _, sk := range skils {
		f := strings.Split(sk, ":")
		result[f[0]] = f[1]
	}
	return result
}

func ResultModifiCharacteristic(ch []Characteristic) []Characteristic {
	for i, char := range ch {
		switch char.Col {
		case 1:
			ch[i].Mod = -5
		case 2:
			ch[i].Mod = -4
		case 3:
			ch[i].Mod = -4
		case 4:
			ch[i].Mod = -3
		case 5:
			ch[i].Mod = -3
		case 6:
			ch[i].Mod = -2
		case 7:
			ch[i].Mod = -2
		case 8:
			ch[i].Mod = -1
		case 9:
			ch[i].Mod = -1
		case 10:
			ch[i].Mod = 0
		case 11:
			ch[i].Mod = 0
		case 12:
			ch[i].Mod = 1
		case 13:
			ch[i].Mod = 1
		case 14:
			ch[i].Mod = 2
		case 15:
			ch[i].Mod = 2
		case 16:
			ch[i].Mod = 3
		case 17:
			ch[i].Mod = 3
		case 18:
			ch[i].Mod = 4
		case 19:
			ch[i].Mod = 4
		case 20:
			ch[i].Mod = 5
		case 21:
			ch[i].Mod = 5
		case 22:
			ch[i].Mod = 6
		case 23:
			ch[i].Mod = 6
		case 24:
			ch[i].Mod = 7
		case 25:
			ch[i].Mod = 7
		case 26:
			ch[i].Mod = 8
		case 27:
			ch[i].Mod = 8
		case 28:
			ch[i].Mod = 9
		case 29:
			ch[i].Mod = 9
		case 30:
			ch[i].Mod = 10
		}
	}
	return ch
}
func ResultExpLvl(ex int64) int {
	if ex < 300 {
		return 1
	} else if ex < 900 {
		return 2
	} else if ex < 2700 {
		return 3
	} else if ex < 6500 {
		return 4
	} else if ex < 14000 {
		return 5
	} else if ex < 23000 {
		return 6
	} else if ex < 34000 {
		return 7
	} else if ex < 48000 {
		return 8
	} else if ex < 64000 {
		return 9
	} else if ex < 85000 {
		return 10
	} else if ex < 100000 {
		return 11
	} else if ex < 120000 {
		return 12
	} else if ex < 140000 {
		return 13
	} else if ex < 165000 {
		return 14
	} else if ex < 195000 {
		return 15
	} else if ex < 225000 {
		return 16
	} else if ex < 265000 {
		return 17
	} else if ex < 305000 {
		return 18
	} else if ex < 355000 {
		return 19
	} else {
		return 20
	}
}
func ResultBonusMaster(lvl int) int {
	if lvl <= 4 {
		return 2
	} else if lvl <= 8 {
		return 3
	} else if lvl <= 12 {
		return 4
	} else if lvl <= 16 {
		return 5
	} else {
		return 6
	}
}

type WeaponT struct {
	Upgrade     int    `json:"upgrade"`
	Description string `json:"description"`
	Type        string `json:"type"` // Например двуручный топор
	Damage      string `json:"damage"`
	UnicBonuses string `json:"unic_bonuses"` //уникальные бонусы например +5 у крону по ограм
}

type Spels struct {
	Name     string `json:"name"`
	Lvl      int    `json:"lvl"`       // Уровень ячейки от 0 до 9
	Damage   string `json:"damage"`    // урон заклинания в стиле 1d6
	TypeSpas string `json:"type_spas"` // Тип спаса Допустим Интелект
	HardSpas int    `json:"hard_spas"` // сложность спаса(должен высчитывать)
}

type Characteristic struct {
	Name string `json:"name"`
	Col  int    `json:"col"`
	Mod  int    `json:"mod,omitempty"`
}

type Char struct {
	NameUser          string            `json:"name_user,omitempty"`    // Имя пользователя
	NameModule        string            `json:"name_module,omitempty"`  // Имя модуля(свободный ввод)
	NameChar          string            `json:"name_char"`              // Имя персоонажа
	Race              string            `json:"race"`                   // Расса
	Characteristic    []Characteristic  `json:"characteristic"`         // Характеристики в формате Название/Значение
	Experience        int64             `json:"experience"`             // Опыт
	Lvl               int               `json:"lvl,omitempty"`          // Уровень перса
	Class             string            `json:"class"`                  // Класс персоонажа
	Weapon            []WeaponT         `json:"weapon"`                 // Всё ОРУЖИЕ перса которым он пользуется
	Skills            []string          `json:"skills"`                 // стандартные скилы по типу атлетики и акробатики, заполнять только взятые\
	BonusMaster       int               `json:"bonus_master,omitempty"` //Бонус мастерства
	UnicueBonusSkills map[string]string `json:"unicue_bonus_skills"`    //Если у когото уесть дополнительно усиление навыка
	Gold              int64             `json:"gold"`                   // Будем учитывать только золото
	Invertar          string            `json:"invertar"`               // Просто вещи писать можно всё что угодно
	Spels             []Spels           `json:"spels"`                  //Заклинания персоонажа
	UnicueSpels       map[string]string `json:"unicue_spels"`           // Уникальные способности в стиле Название:Эффект
	Resurses          string            `json:"resurses"`               //уникальный ресурс персоонажа(если есть) в идеале заполнять Название + кол-во
	Description       string            `json:"description"`            // Описание персоонажа
	ImageUrl          string            `json:"image_url"`              // Картинка перса
	NumModule         int               `json:"num_module"`
}
