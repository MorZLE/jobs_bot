package service

import (
	"github.com/MorZLE/jobs_bot/model"
)

type Service interface {
	SaveResume(user model.Student) error
	Get(id int64) (model.Student, error)
	Delete(id int64, category string) error
	GetResume(category string, count int, direction string, wantStatus string) (model.Student, int, error)

	BanUser(idx int, category string) error
	PublishUser(idx int, category string) error
	DeclineUser(idx int, category string) error
	Statistics() (string, error)
	UnbanUser(user, flag string) error
	ViewBanList() (string, error)
	NewAdmin(username string) (string, error)
	AuthNewAdmin(id int64, username, url string) error
	GetAdmins() []int64
}
