package db

type User struct {
	Uuid  string	`gorm:"primaryKey"`
	Phone string	`gorm:"unique"`
	Name  string
}
