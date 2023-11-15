package service

import (
	"github.com/MorZLE/jobs_bot/model"
)

type Service interface {
	SaveResume(user model.Student) error
	Get(id int64) (model.Student, error)
	Delete(id int64) error
	GetResume(category string, count int) (model.Student, error)
}
