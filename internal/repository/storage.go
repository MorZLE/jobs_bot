package repository

import "github.com/MorZLE/jobs_bot/model"

type Storage interface {
	Set(student model.Student) error
	Get(id int64) (model.Student, error)
	Delete(id int64, category string) error
	Close()
	GetOneResume(category string, direction string, count int, wantStatus string) (model.Student, error)

	BanUser(idx int, category string) error
	PublishUser(idx int, category string) error
	DeclineUser(idx int, category string) error
	Statistics() (map[string][]model.Student, error)
	UnbanUsername(username string) error
	UnbanTgID(tgid int64) error
	ViewBanList() ([]model.BanUser, error)
	NewAdminURL(username, url string) error
	CheckUrlAdmin(username, url string) error
	CreateAdmin(username string, id int64) error
	GetAdmins() ([]model.Admin, error)
	DeleteUrlInvaite(username, url string) error
	DeleteAdmin(username string) error
}
