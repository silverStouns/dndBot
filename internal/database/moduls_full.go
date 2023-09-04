package database

import (
	"database/sql"
	"dndBot/internal/pkg/logger"
	"errors"
)

// Вытягивает всю инфу по всем модулям(игроков нужно стягивать оттдельно)
func GetModulsInfo(conn *DBConnector) (err error) {
	var (
		moduleName  sql.NullString
		masterName  sql.NullString
		description sql.NullString
	)
	rows, err := conn.Connector.Query(`select module_name, m.name, description
    from moduls
    inner join main.user u on u.id = moduls.gamers
    inner join main.user m on m.id = moduls.master_id
group by moduls.module_name;`)
	for rows.Next() {
		err = rows.Scan(&moduleName, &masterName, &description)
	}

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil
	case err != nil:
		logger.Debug("GetModulsInfo err")
		return err
	default:
		return nil
	}
}
