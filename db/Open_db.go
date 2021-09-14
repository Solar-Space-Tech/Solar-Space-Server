package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// db 对象可以 “ .Close() ”
func Open_db() *sql.DB {
	db_, err := os.Open("./db.json")
	if err != nil {
		log.Panicln(err)
	}
	var (
		db_sc struct {
			DB_type string `json:"db_type"`
			Host    string `json:"host"`
		}
	)
	if err := json.NewDecoder(db_).Decode(&db_sc); err != nil {
		log.Panicln(err)
	}
	db, err := sql.Open(db_sc.DB_type, db_sc.Host)
	if err != nil {
		log.Panicln(err)
		fmt.Println(err)
	}
	db.SetConnMaxLifetime(100)
	db.SetMaxIdleConns(5)

	return db
}
