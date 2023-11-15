package model

import "github.com/jinzhu/gorm"

type Student struct {
	gorm.Model
	Id       int    `gorm:"primary_key" gorm:"AUTO_INCREMENT"`
	Tgid     int64  `gorm:"not null" gorm:"unique"`
	Fio      string `gorm:"not null"`
	Group    string `gorm:"not null"`
	Resume   string `gorm:"not null"`
	Category string `gorm:"not null"`
	Status   string `gorm:"not null"`
}
type Employee struct {
	gorm.Model
	Id      int    `gorm:"primary_key" gorm:"AUTO_INCREMENT"`
	Company string `gorm:"not null"`
	Resume  string `gorm:"not null"`
}

type Status struct {
	Id           int    `gorm:"primary_key" gorm:"AUTO_INCREMENT"`
	Registration string `gorm:"not null"`
	Moderation   string `gorm:"not null"`
	Active       string `gorm:"not null"`
	Disabled     string `gorm:"not null"`
	Ban          string `gorm:"not null"`
}
