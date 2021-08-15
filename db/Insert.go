package db

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func Insert_mixin(phone, uuid, name string) bool  {
    db := Open_db()
    if err := db.Ping(); err != nil {
        log.Panicln(err)
        return false
    }
    defer db.Close()

    stmt, err := db.Prepare("INSERT usermixin SET phone=?,uuid=?,name=?")
    checkErr(err)

    res, err := stmt.Exec(phone, uuid, name)
    checkErr(err)

    id, err := res.LastInsertId()
    checkErr(err)

    fmt.Println(id)

    return  true
}


func checkErr(err error) {
    if err != nil {
        log.Panicln(err)
    }
}