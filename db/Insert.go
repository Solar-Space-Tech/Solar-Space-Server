package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func Insert_mixin(phone, uuid string) bool  {
    db := Open_db()
    if err := db.Ping(); err != nil {
        log.Panicln(err)
        return false
    }
    stmt, err := db.Prepare("INSERT usermixin SET phone=?,uuid=?")
    checkErr(err)

    res, err := stmt.Exec(phone, uuid)
    checkErr(err)

    id, err := res.LastInsertId()
    checkErr(err)

    fmt.Println(id)

    db.Close()

    return  true
}


func checkErr(err error) {
    if err != nil {
        log.Panicln(err)
    }
}