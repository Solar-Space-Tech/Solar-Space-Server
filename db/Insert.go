package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func Insert_mixin(db *sql.DB, phone, uuid string)  {
    stmt, err := db.Prepare("INSERT usermixin SET phone=?,uuid=?")
    checkErr(err)

    res, err := stmt.Exec(phone, uuid)
    checkErr(err)

    id, err := res.LastInsertId()
    checkErr(err)

    fmt.Println(id)
}


func checkErr(err error) {
    if err != nil {
        log.Panicln(err)
    }
}