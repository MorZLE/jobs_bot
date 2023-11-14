package service

import (
	"github.com/MorZLE/jobs_bot/model"
	"gopkg.in/telebot.v3"
)

type Service interface {
	SaveResume(id int64, doc telebot.File, user model.Student) error
	Get(id int64) (model.Student, error)
}
