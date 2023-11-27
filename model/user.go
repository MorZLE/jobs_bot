package model

import "gorm.io/gorm"

type User struct {
	Student             Student
	Type                string //Тип пользователя
	Nqest               int    //Номер вопроса при регистрации
	EmployeeCount       int    //Номер резюме в очереди при просмотре
	EmployeeCategory    string //Категория резюме при просмотре
	EmployeeSetCategory bool   //Флаг смены категории
}

type AdminInvait struct {
	gorm.Model
	Username string
	Url      string
}

type Admin struct {
	Username string
	Lvl      int
	Tgid     int64
}
