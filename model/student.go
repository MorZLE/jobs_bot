package model

import "github.com/jinzhu/gorm"

type Student struct {
	gorm.Model
	Id       int    `gorm:"primary_key" gorm:"AUTO_INCREMENT"`
	Tgid     int64  `gorm:"not null" gorm:"unique"`
	Username string `gorm:"not null"`
	Fio      string `gorm:"not null"`
	Group    string `gorm:"not null"`
	Resume   string `gorm:"not null"`
	Category string `gorm:"not null"`
	Status   string `gorm:"not null"`
}
type BanUser struct {
	gorm.Model
	Tgid     int64  `gorm:"not null" gorm:"unique"`
	Username string `gorm:"not null"`
}
