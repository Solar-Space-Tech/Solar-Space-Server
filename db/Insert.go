package db

import (
	"github.com/fox-one/mixin-sdk-go"
	_ "github.com/go-sql-driver/mysql"
)

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

func Insert_subWallet(s *mixin.Keystore) bool {
	db, _ := Open_db()

	defer db.Close()

	wallet_info := SubWallet{
		ClientID:   s.ClientID,
		SessionID:  s.SessionID,
		PrivateKey: s.PrivateKey,
		PinToken:   s.PinToken,
	}
	db.Create(&wallet_info)

	return true
}
