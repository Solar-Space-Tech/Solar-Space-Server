package db

type User struct {
	Uuid  string	`gorm:"primaryKey"`
	Phone string	`gorm:"unique"`
	Name  string
}

type SubWallet struct {
	ClientID   string `json:"client_id"`
	SessionID  string `json:"session_id"`
	PrivateKey string `gorm:"type:TEXT"`
	PinToken   string `json:"pin_token"`
}
