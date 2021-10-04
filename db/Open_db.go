package db

import (
	"encoding/json"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// db 对象可以 “ .Close() ”
func Open_db() (*gorm.DB, error) {
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
	db, err := gorm.Open(db_sc.DB_type, db_sc.Host)
	if err != nil {
		return nil, err
	}
	checkErr(db.DB().Ping())

	db.AutoMigrate(&User{}) // Generate sheet by struct

	return db, nil
}

func checkErr(err error) {
	if err != nil {
		log.Panicln(err)
	}
}
