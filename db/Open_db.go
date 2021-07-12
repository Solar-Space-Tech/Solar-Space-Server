package db

import (
    "database/sql"
    "fmt"
	"log"

    _ "github.com/go-sql-driver/mysql"
)
// DB 对象可以 “ .Close() ”
func Open_db() *sql.DB {
	db, err := sql.Open("mysql", "sst_server:sexy0756@tcp(rm-bp1y56w4272v2u5frzo.mysql.rds.aliyuncs.com:3306)/sst?charset=utf8")
    if err != nil {
		log.Panicln(err)
		fmt.Println(err)
	} 
	db.SetConnMaxLifetime(100)
	db.SetMaxIdleConns(5)

	return db;
}