package service

import (
	"errors"
	"fmt"
	"github.com/MorZLE/jobs_bot/config"
	"github.com/MorZLE/jobs_bot/constants"
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

func (s *serviceImpl) Delete(id int64, category string) error {
	if err := s.db.Delete(id, category); err != nil {
		return err
	}
	err := os.Remove(fmt.Sprintf("%s\\src\\resume\\%d.pdf", s.dir, id))
	if err != nil {
		return err
	}
	return nil
}

func (s *serviceImpl) GetResume(category string, count int, direction string) (model.Student, int, error) {
	user, err := s.db.GetOneResume(category, direction, count)
	if err != nil {
		if errors.Is(err, constants.ErrDeleteResume) {
			switch direction {
			case constants.Next:
				count++
				user, c, err := s.GetResume(category, count, direction)
				return user, c, err
			case constants.Prev:
				if count > 1 {
					count--
				} else {
					return model.Student{}, 0, constants.ErrNotFound
				}
				user, c, err := s.GetResume(category, count, direction)
				return user, c, err
			}
			return model.Student{}, 0, err
		}
		return model.Student{}, count, err
	}
	return user, count, nil
}
