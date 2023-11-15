package repository

import "github.com/MorZLE/jobs_bot/model"

type Storage interface {
	Set(student model.Student) error
	Get(id int64) (model.Student, error)
	Delete(id int64) error
	Close()
	GetOneResume(category string, count int) (model.Student, error)
}
