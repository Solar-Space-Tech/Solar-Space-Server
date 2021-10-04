package db

import (
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Uuid  string
	Phone string
	Name  string
}

func Insert_mixin(phone, uuid, name string) bool {
	db, _ := Open_db()

	defer db.Close()

	user_info := User{
		Uuid:  uuid,
		Phone: phone,
		Name:  name,
	}
	db.Create(&user_info)

	return true
}
