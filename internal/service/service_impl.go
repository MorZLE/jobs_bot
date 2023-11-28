package service

import (
	"errors"
	"fmt"
	"github.com/MorZLE/jobs_bot/config"
	"github.com/MorZLE/jobs_bot/constants"
	"github.com/MorZLE/jobs_bot/internal/repository"
	"github.com/MorZLE/jobs_bot/logger"
	"github.com/MorZLE/jobs_bot/model"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/speps/go-hashids"
	"os"
	"strconv"
	"time"
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
	pdfPath := fmt.Sprintf("%s\\src\\resume\\%s", s.dir, user.Resume)
	err := api.ValidateFile(pdfPath, nil)
	if err != nil {
		logger.Error("Failed to open PDF file:", err)
		return constants.ErrOpenFile
	}
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

func (s *serviceImpl) GetResume(category string, count int, direction string, wantStatus string) (model.Student, int, error) {
	user, err := s.db.GetOneResume(category, direction, count, wantStatus)
	if err != nil {
		if errors.Is(err, constants.ErrLastResume) {
			return user, count, err
		}
		if errors.Is(err, constants.ErrDeleteResume) {
			switch direction {
			case constants.Next:
				count++
				return s.GetResume(category, count, direction, wantStatus)
			case constants.Prev:
				if count >= 1 {
					count--
				} else {
					return model.Student{}, 0, constants.ErrNotFound
				}
				return s.GetResume(category, count, direction, wantStatus)
			}
			return model.Student{}, 0, err
		}
		return model.Student{}, count, err
	}
	return user, count, nil
}
func (s *serviceImpl) BanUser(idx int, category string) error {
	err := s.Delete(int64(idx), category)
	if err != nil {
		return err
	}
	return s.db.BanUser(idx, category)
}
func (s *serviceImpl) PublishUser(idx int, category string) error {
	return s.db.PublishUser(idx, category)
}
func (s *serviceImpl) DeclineUser(idx int, category string) error {
	return s.db.DeclineUser(idx, category)
}

func (s *serviceImpl) Statistics() (string, error) {
	m, _ := s.db.Statistics()
	res := "Статистика:\n"
	for k, v := range m {
		res += fmt.Sprintf("%s: %d\n", k, len(v))
	}
	return res, nil
}

func (s *serviceImpl) UnbanUser(user, flag string) error {
	switch flag {
	case constants.Username:
		return s.db.UnbanUsername(user)
	case constants.TgID:
		tgid, err := strconv.ParseInt(user, 10, 64)
		if err != nil {
			return err
		}
		return s.db.UnbanTgID(tgid)
	}
	return nil
}

func (s *serviceImpl) ViewBanList() (string, error) {
	res := "Список забаненых пользователей:\n"
	userban, err := s.db.ViewBanList()
	if err != nil {
		return "", err
	}
	for _, v := range userban {

		res += fmt.Sprintf("%s: %d\n", v.Username, int(v.Tgid))
	}
	return res, nil
}

func (s *serviceImpl) NewAdmin(username string) (string, error) {
	hd := hashids.NewData()
	hd.MinLength = 6
	h, err := hashids.NewWithData(hd)
	if err != nil {
		logger.Error("Error NewWithData:", err)
		return "", err
	}
	url, err := h.Encode([]int{time.Now().Nanosecond()})
	if err != nil {
		logger.Error("Error Encode:", err)
		return "", err
	}
	err = s.db.NewAdminURL(username, url)
	return url, nil
}

func (s *serviceImpl) AuthNewAdmin(id int64, username, url string) error {
	if err := s.db.CheckUrlAdmin(username, url); err != nil {
		return err
	}
	if err := s.db.CreateAdmin(username, id); err != nil {
		return err
	}
	if err := s.db.DeleteUrlInvaite(username, url); err != nil {
		return err
	}
	return nil
}

func (s *serviceImpl) GetAdmins() (string, error) {
	res, err := s.db.GetAdmins()
	if err != nil {
		return "", err
	}
	adm := "Список администраторов:\n"
	for _, v := range res {
		adm += fmt.Sprintf("@%s\n", v.Username)
	}
	return adm, nil
}

func (s *serviceImpl) DeleteAdmin(username string) error {
	return s.db.DeleteAdmin(username)
}
