package db

import (
	"fmt"

	"github.com/fox-one/pkg/property"
	_ "github.com/go-sql-driver/mysql"
)

func Query_uuid_by_phone(phone string) string {
	db, err := Open_db()
	checkErr(err)

	defer db.Close()

	var qUser User
	d := db.Where("phone = ?", phone).First(&qUser)
	v, ok := d.Value.(*User)
	if ok {
		return v.Uuid
	}
	return "The phone is not exist."
}

func If_old_user(uuid, phone string) bool {
	db, err := Open_db()
	checkErr(err)
	defer db.Close()

	var qUsers []*User
	u := db.Where("uuid = ?", uuid).Or("phone = ?", phone).Find(&qUsers)
	rows, ok := u.Value.(*[]*User)
	if ok {
		fmt.Println("lalala")
		fmt.Printf("v: %v\n", rows)

		for _, row := range *rows {
			if uuid == row.Uuid || phone == row.Phone {
				return true
			}
		}
	}
	return false
}

// Get offset
func Get_utxo(key string) (property.Value, error) {
	db, err := Sqlite_open_db()
	checkErr(err)
	defer db.Close()

	var p Property
	// err := s.db.Where(tableName+".key = ?", key).First(&p).Error
	err = db.Where("key = ?", key).First(&p).Error
	if IsErrorNotFound(err) {
		err = nil
	}
	return p.Value, err
}