package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		fmt.Println(err)
		return err
	}

	var schema = `CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(64) DEFAULT "",
    comment TEXT DEFAULT "",
    repeat VARCHAR(128) DEFAULT "");

	CREATE INDEX dateindex ON scheduler(date);
    `

	if install {
		_, err := DB.Exec(schema)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return err
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
