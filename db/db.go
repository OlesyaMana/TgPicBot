package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/glebarez/go-sqlite"
)

const (
	dbFileName = "bot_db.db"
	DogsColumn = "receive_dogs"
)

func StartDB() (*sql.DB, error) {
	_, err := os.Stat(dbFileName)
	dbNotExists := os.IsNotExist(err)
	if dbNotExists {
		_, err = os.Create(dbFileName)
		if err != nil {
			return nil, err
		}
	}

	db, err := sql.Open("sqlite", dbFileName)
	if err != nil {
		return nil, err
	}

	if dbNotExists {
		_, err = db.Exec("CREATE TABLE Bot_Users (id BIGINT, login TEXT, receive_dogs INTEGER)")
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

func AddNewUser(db *sql.DB, id int64, login string) error {
	_, err := db.Exec(fmt.Sprintf("INSERT INTO Bot_Users VALUES (%d, '%s', 0)", id, login))
	return err
}

func DeleteUser(db *sql.DB, id int64) error {
	_, err := db.Exec(fmt.Sprintf("DELETE FROM Bot_Users WHERE id=%d", id))
	return err
}

func UpdateUser(db *sql.DB, userId int64, columnName string, value int) error {
	var query string
	switch columnName {
	case DogsColumn:
		query = fmt.Sprintf("UPDATE Bot_Users SET %s = %d WHERE id=%d", DogsColumn, value, userId)

	default:
		return fmt.Errorf("unknown column name")
	}

	_, err := db.Exec(query)
	return err
}

func GetBotUsers(db *sql.DB) (map[int64]BotUser, error) {
	rows, err := db.Query("SELECT * FROM Bot_Users")
	if err != nil {
		return nil, err
	}

	forReturn := make(map[int64]BotUser)
	for rows.Next() {
		var user BotUser
		err := rows.Scan(&user.Id, &user.Login, &user.ReceiveDogs)
		if err != nil {
			return nil, err
		}

		forReturn[user.Id] = user
	}

	return forReturn, nil
}

type BotUser struct {
	Id          int64  `db:"id"`
	Login       string `db:"login"`
	ReceiveDogs int    `db:"receive_dogs"`
}
