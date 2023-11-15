package service

import (
	"fmt"
	"github.com/MorZLE/jobs_bot/config"
	"github.com/MorZLE/jobs_bot/model"
	"github.com/MorZLE/jobs_bot/repository"
	"os"
)

func NewService(cnf *config.Config, db repository.Storage) Service {

	return &serviceImpl{
		db:  db,
		dir: cnf.Dir,
	}
}

type serviceImpl struct {
	db  repository.Storage
	dir string
}

func (s *serviceImpl) SaveResume(user model.Student) error {
	if err := s.db.Set(user); err != nil {
		return err
	}
	return nil
}

func (s *serviceImpl) Get(id int64) (model.Student, error) {
	user, err := s.db.Get(id)
	if err != nil {
		return model.Student{}, err
	}
	return user, nil
}

func (s *serviceImpl) Delete(id int64) error {
	if err := s.db.Delete(id); err != nil {
		return err
	}
	err := os.Remove(fmt.Sprintf("%s\\src\\resume\\%d.pdf", s.dir, id))
	if err != nil {
		return err
	}
	return nil
}
