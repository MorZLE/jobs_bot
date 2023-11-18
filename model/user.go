package model

type User struct {
	Student             Student
	Type                string //Тип пользователя
	Nqest               int    //Номер вопроса при регистрации
	EmployeeCount       int    //Номер резюме в очереди при просмотре
	EmployeeCategory    string //Категория резюме при просмотре
	EmployeeSetCategory bool   //Флаг смены категории
}
