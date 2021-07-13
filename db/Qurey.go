package db

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func Query_uuid_by_phone(phone string) string {
    db := Open_db()
    if err:= db.Ping(); err != nil {
        log.Panicln(err)
    }
    defer db.Close()

	rows, err := db.Query("SELECT phone, uuid FROM usermixin")
    checkErr(err)
	
	for rows.Next() {
        var phone_number string
        var uuid string
        err = rows.Scan(&phone_number, &uuid)
        checkErr(err)
        if phone_number == phone{
            return uuid
        }
    }
    defer rows.Close()

    return "The phone is not exist."
}
