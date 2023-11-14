package service

import (
	"fmt"
	"github.com/MorZLE/jobs_bot/model"
	"github.com/MorZLE/jobs_bot/repository"
	"gopkg.in/telebot.v3"
	"io"
	"os"
)

func NewService(db repository.Storage) Service {
	return &serviceImpl{
		db: db,
	}
}

type serviceImpl struct {
	db repository.Storage
}

func (s *serviceImpl) SaveResume(id int64, doc telebot.File, user model.Student) error {
	if err := s.db.Set(user); err != nil {
		return err
	}
	pdfPath := fmt.Sprintf("../src/pdf/%d.pdf", id)
	newFile, err := os.Create(pdfPath)
	if err != nil {
		return err
	}
	defer newFile.Close()

	// Копируем содержимое оригинального файла в новый файл
	_, err = io.Copy(newFile, doc.FileReader)
	if err != nil {
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

func (s *serviceImpl) Delete() {

}
