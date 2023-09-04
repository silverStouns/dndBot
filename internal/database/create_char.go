package database

import (
	"database/sql"
	"dndBot/internal/pkg/logger"
	"errors"
	"fmt"
)

func GetGamerId(conn *DBConnector, tgId string) (id int64, err error) {
	err = conn.Connector.
		QueryRow(`select id
		from user
		where telegram_id = $1;`, tgId).Scan(&id)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return 0, nil
	case err != nil:
		logger.Debug("getGamerId")
		return 0, err
	default:
		logger.Debug("getGamerId")
		return id, nil
	}
}

func CreateChar(conn *DBConnector, char Char, gamerId int64) error {
	logger.Debug("char insert")
	_, err := conn.Connector.Exec(`insert into chars(gamer_id, module_id, name,
                  characteristic, exp, class, weapon, skills, 
                  unic_skills, gold, invertory, spels,
                  unic_spels, unic_resurses, desription, race)
values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)`, gamerId, char.NumModule, char.NameChar, converToCharcs(char.Characteristic), char.Experience, char.Class, converToWeapon(char.Weapon), converToSkills(char.Skills), converToUnicSkills(char.UnicueBonusSkills),
		char.Gold, char.Invertar, converToSpels(char.Spels), converToUnicSpels(char.UnicueSpels), char.Resurses, char.Description, char.Race)
	return err
}

func converToCharcs(characteristic []Characteristic) (charcs string) {
	//[{"name":"Сила","col":21},
	//{"name":"Ловкость","col":16},
	//{"name":"Телосложение","col":20},
	//{"name":"Интеллект","col":14},
	//{"name":"Мудрость","col":14},
	//{"name":"Харизма","col":14}]
	v := ""
	iter := 0
	for _, j := range characteristic {
		if iter < len(characteristic)-1 {
			jj := fmt.Sprintf(`{"name":"%v","col":%v},`, j.Name, j.Col)
			v += jj
		} else {
			jj := fmt.Sprintf(`{"name":"%v","col":%v}`, j.Name, j.Col)
			v += jj
		}
		iter++
	}
	return fmt.Sprintf("[%s]", v)
}

func converToWeapon(weapon []WeaponT) (charcs string) {
	//[{"upgrade":1,"description":"Богорез","type":"Одноручная коса","damage":"1d12","unic_bonuses":"Способна ослепить всех в радиусе 15фт, сложность спаса 14"},
	//	{"upgrade":0,"description":"Здоровенный сука молот!","type":"Двуручный молот","damage":"1d8(Димас пидорас)","unic_bonuses":""}]
	v := ""
	iter := 0
	for _, j := range weapon {
		if iter < len(weapon)-1 {
			jj := fmt.Sprintf(`{"upgrade":%v,"description":"%v","type":"%v","damage":"%v","unic_bonuses":"%v"},`, j.Upgrade, j.Description, j.Type, j.Damage, j.UnicBonuses)
			v += jj
		} else {
			jj := fmt.Sprintf(`{"upgrade":%v,"description":"%v","type":"%v","damage":"%v","unic_bonuses":"%v"}`, j.Upgrade, j.Description, j.Type, j.Damage, j.UnicBonuses)
			v += jj
		}
		iter++
	}
	return fmt.Sprintf("[%s]", v)
}

func converToSkills(weapon []string) (charcs string) {
	//Медицина, Атлетика, Наблюдательность

	for i, j := range weapon {
		if i != 0 {
			charcs += fmt.Sprintf(", %v", j)
		} else {
			charcs += fmt.Sprintf(j)
		}
	}
	return
}

func converToUnicSkills(weapon map[string]string) (charcs string) {
	//Харизма:1
	iter := 0
	for i, j := range weapon {
		if iter < len(weapon)-1 {
			charcs += fmt.Sprintf("%v:%v,", i, j)
		} else {
			charcs += fmt.Sprintf("%v:%v", i, j)
		}
		iter++
	}
	return
}

func converToSpels(spels []Spels) (charcs string) {
	//[{"lvl":0,"name":"Сотворение костра","damage":"1d8","type_spas":"ловкость","hard_spas":15},
	//{"lvl":0,"name":"Электрическая плеть","damage":"1d8","type_spas":"сила","hard_spas":15}]
	v := ""
	iter := 0
	for _, j := range spels {
		if iter < len(spels)-1 {
			jj := fmt.Sprintf(`{"lvl":%v,"name":"%v","damage":"%v","type_spas":"%v","hard_spas":%v},`, j.Lvl, j.Name, j.Damage, j.TypeSpas, j.HardSpas)
			v += jj
		} else {
			jj := fmt.Sprintf(`{"lvl":%v,"name":"%v","damage":"%v","type_spas":"%v","hard_spas":%v}`, j.Lvl, j.Name, j.Damage, j.TypeSpas, j.HardSpas)
			v += jj
		}
		iter++
	}
	return fmt.Sprintf("[%s]", v)
}

func converToUnicSpels(weapon map[string]string) (charcs string) {
	iter := 0
	//Гарпун:Выстрел из гарпуна на 30фт:::Скрытый щит:В руке спрятан щит на 4 армора:::Форсаж:3 заряда, тот же эффект что у заклинания ускорение
	for i, j := range weapon {
		if iter != 0 {
			charcs += fmt.Sprintf(":::%v:%v", i, j)
		} else {
			charcs += fmt.Sprintf("%v:%v", i, j)
		}
		iter++
	}
	return
}
