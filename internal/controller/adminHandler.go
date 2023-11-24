package controller

import (
	"errors"
	"fmt"
	"github.com/MorZLE/jobs_bot/constants"
	"github.com/MorZLE/jobs_bot/logger"
	bot "gopkg.in/telebot.v3"
	"log"
)

var (
	menuAdmin          = &bot.ReplyMarkup{}
	btnViewResumeAdmin = menuAdmin.Data("Модерация резюме", "modResume")
	btnBanUser         = selector.Data("Забанить пользователя", "banUser")
	btnDeclineUser     = selector.Data("Отклонить резюме", "declineUser")
	btnPublishUser     = selector.Data("Опубликовать резюме", "publishUser")
)

func (h *Handler) CheckAdmin(c bot.Context) error {
	id := c.Sender().ID
	if id != h.admin {
		h.bot.Send(c.Chat(), "Вы не являетесь администратором")
		return nil
	}
	h.mutex.Lock()
	user := h.user[id]
	user.Type = constants.Admin
	h.user[id] = user
	h.mutex.Unlock()
	return nil
}

func (h *Handler) AdminMenu(c bot.Context) error {
	menuAdmin.Inline(
		menuAdmin.Row(
			btnViewResumeAdmin,
		),
	)
	_, err := h.bot.Send(c.Chat(), "Админ меню", menuAdmin)
	if err != nil {
		logger.Error("ошибка AdminMenu", err)
	}
	return nil
}

func (h *Handler) btnViewResumeAdmin(c bot.Context) error {
	h.CheckAdmin(c)
	h.btnCategorySelect(c)
	return nil
}

func (h *Handler) ModerationAdmin(c bot.Context, dir string) error {
	h.CheckAdmin(c)
	id := c.Sender().ID
	if id == h.admin {
		mUser := h.user[id]
		user, count, err := h.s.GetResume(mUser.EmployeeCategory, mUser.EmployeeCount, dir, constants.StatusModeration)
		if count > 0 {
			selector.Inline(
				selector.Row(btnPrev, btnNext),
				selector.Row(btnPublishUser),
				selector.Row(btnDeclineUser),
				selector.Row(btnBanUser),
			)
		} else if count == 0 {
			selector.Inline(
				selector.Row(btnLock, btnNext),
				selector.Row(btnPublishUser),
				selector.Row(btnDeclineUser),
				selector.Row(btnBanUser),
			)
		}
		if err != nil {
			if errors.Is(err, constants.ErrNotCategory) {
				h.bot.Send(c.Chat(), "Категория не найдена")
				return nil
			}
			if errors.Is(err, constants.ErrNotFound) {
				h.bot.Send(c.Chat(), "Резюме закончились")

				return nil
			}
			if errors.Is(err, constants.ErrNotResume) {
				h.bot.Send(c.Chat(), "Нету резюме в данной категории")
				return nil
			}
			if errors.Is(err, constants.ErrLastResume) {
				selector.Inline(
					selector.Row(btnPrev, btnLock),
					selector.Row(btnPublishUser),
					selector.Row(btnDeclineUser),
					selector.Row(btnBanUser),
				)
				if count == 0 {
					selector.Inline(
						selector.Row(btnPublishUser),
						selector.Row(btnDeclineUser),
						selector.Row(btnBanUser),
					)
				}
			} else {
				h.bot.Send(c.Chat(), "Что то пошло не так")
				log.Println(err)
			}
		}

		mUser.EmployeeCount = count
		urlPDF := fmt.Sprintf("src\\resume\\%s", user.Resume)
		resume := fmt.Sprintf("ФИО: %s\nГруппа: %s\nКатегория: %s\nПользователь: @%s", user.Fio, user.Group, user.Category, user.Username)
		file := &bot.Photo{
			File:    bot.FromDisk(urlPDF),
			Caption: resume,
		}
		if mUser.EmployeeSetCategory {
			_, err = h.bot.Send(c.Chat(), file, selector)
			if err != nil {
				h.ViewResumeStudents(c, constants.Next)
				logger.Error("ошибка загрузки файла в резюме", err)
			}
		} else {
			_, err = h.bot.Edit(c.Message(), file, selector)
			if err != nil {
				h.ViewResumeStudents(c, constants.Next)
				logger.Error("ошибка загрузки файла в резюме", err)
			}
		}
		mUser.EmployeeSetCategory = false
		h.mutex.Lock()
		h.user[id] = mUser
		h.mutex.Unlock()
		return nil
	}
	return nil
}

func (h *Handler) btnBanUser(c bot.Context) error {
	id := c.Sender().ID
	mUser := h.user[id]
	err := h.s.BanUser(mUser.EmployeeCount, mUser.EmployeeCategory)
	if err != nil {
		h.bot.Send(c.Chat(), "Что то пошло не так")
		log.Println(err)
	}
	return nil
}

func (h *Handler) btnDeclineUser(c bot.Context) error {
	id := c.Sender().ID
	mUser := h.user[id]
	err := h.s.DeclineUser(mUser.EmployeeCount, mUser.EmployeeCategory)
	if err != nil {
		h.bot.Send(c.Chat(), "Что то пошло не так")
		log.Println(err)
	}
	return nil
}
func (h *Handler) btnPublishUser(c bot.Context) error {
	id := c.Sender().ID
	mUser := h.user[id]
	err := h.s.PublishUser(mUser.EmployeeCount, mUser.EmployeeCategory)
	if err != nil {
		h.bot.Send(c.Chat(), "Что то пошло не так")
		log.Println(err)
	}
	return nil
}
