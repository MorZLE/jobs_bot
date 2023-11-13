package service

import "github.com/MorZLE/jobs_bot/repository"

func NewService(db repository.Storage) Service {
	return &serviceImpl{
		db: db,
	}
}

type serviceImpl struct {
	db repository.Storage
}

func (s *serviceImpl) Set() {

}

func (s *serviceImpl) Get() {

}
func (s *serviceImpl) Delete() {

}
