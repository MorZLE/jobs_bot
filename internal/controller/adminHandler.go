package controller

import (
	"errors"
	"fmt"
	"github.com/MorZLE/jobs_bot/constants"
	"github.com/MorZLE/jobs_bot/logger"
	bot "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
	"log"
)

var (
	menuAdmin          = &bot.ReplyMarkup{}
	btnViewResumeAdmin = menuAdmin.Data("Модерация резюме", "modResume")
	btnBanUser         = selector.Data("Забанить пользователя", "banUser")
	btnAdminCommand    = selector.Data("Команды администратора", "btnAdminCommand")
	btnStatistics      = selector.Data("Статистика", "statistics")
	btnViewBanList     = selector.Data("Список забаненных", "viewBanList")
	btnDeclineUser     = selector.Data("Отклонить резюме", "declineUser")
	btnPublishUser     = selector.Data("Опубликовать резюме", "publishUser")
)

func (h *Handler) AdminHandler() {
	adminOnly := h.bot.Group()
	adminOnly.Use(middleware.Whitelist(h.admins...))
	adminOnly.Handle(&adminMenu, h.AdminMenu)
	adminOnly.Handle(&btnViewResumeAdmin, h.btnViewResumeAdmin)
	adminOnly.Handle(&btnBanUser, h.btnBanUser)
	adminOnly.Handle(&btnDeclineUser, h.btnDeclineUser)
	adminOnly.Handle(&btnPublishUser, h.btnPublishUser)
	adminOnly.Handle(&btnAdminCommand, h.btnAdminCommand)
	adminOnly.Handle(&btnStatistics, h.btnStatistics)
	adminOnly.Handle(&btnViewBanList, h.btnViewBanList)

	adminOnly.Handle("/unbanu", h.unbanUsername) //команда разбана по username
	adminOnly.Handle("/unbanid", h.unbanID)      //команда разбана по TGID
	adminOnly.Handle("/newadmin", h.newadmin)    //команда создания админа
}

func (h *Handler) AuthNewAdmin(c bot.Context) error {
	url := c.Args()[0]
	id := c.Sender().ID
	username := c.Sender().Username
	err := h.s.AuthNewAdmin(id, username, url)
	if err != nil {
		logger.Error("ошибка AuthNewAdmin", err)
		return err
	}
	h.admins = append(h.admins, id)
	h.bot.Send(c.Chat(), "Вы авторизованы как администратор")
	h.CheckAdmin(c)
	h.HandlerStart(c)
	return nil
}

func (h *Handler) CheckAdmin(c bot.Context) error {
	id := c.Sender().ID
	h.mutex.Lock()
	user := h.user[id]
	user.Type = constants.Admin
	h.user[id] = user
	h.mutex.Unlock()
	return nil
}

func (h *Handler) AdminMenu(c bot.Context) error {
	menuAdmin.Inline(
		menuAdmin.Row(btnViewResumeAdmin),
		menuAdmin.Row(btnAdminCommand),
		menuAdmin.Row(btnStatistics),
		menuAdmin.Row(btnViewBanList),
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
			h.bot.Send(c.Chat(), "Нету новых резюме в данной категории")
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
			logger.Error("ошибка при модерации резюме", err)
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

func (h *Handler) btnBanUser(c bot.Context) error {
	id := c.Sender().ID
	mUser := h.user[id]
	err := h.s.BanUser(mUser.EmployeeCount, mUser.EmployeeCategory)
	if err != nil {
		h.bot.Send(c.Chat(), "Что то пошло не так")
		logger.Error("ошибка при бане", err)
	}
	return nil
}
func (h *Handler) btnAdminCommand(c bot.Context) error {
	res := "Команды администратора:\n"
	res += "/unbanu имя\n"
	res += "/unbanid TgID\n"
	h.bot.Send(c.Chat(), res)
	return nil
}

func (h *Handler) btnStatistics(c bot.Context) error {
	statistic, err := h.s.Statistics()
	if err != nil {
		h.bot.Send(c.Chat(), "Что то пошло не так")
		logger.Error("ошибка при получении статистики", err)
	}
	_, err = h.bot.Send(c.Chat(), statistic)
	if err != nil {
		h.bot.Send(c.Chat(), "Что то пошло не так")
		logger.Error("ошибка при отправке статистики", err)
	}
	return nil
}
func (h *Handler) unbanID(c bot.Context) error {
	arg := c.Args()[0]
	err := h.s.UnbanUser(arg, constants.TgID)
	if err != nil {
		h.bot.Send(c.Chat(), "Что то пошло не так")
		logger.Error("ошибка при unbanId", err)
	}
	h.bot.Send(c.Chat(), "Пользователь разбанен")
	return nil
}

func (h *Handler) unbanUsername(c bot.Context) error {
	arg := c.Args()[0]
	err := h.s.UnbanUser(arg, constants.Username)
	if err != nil {
		if errors.Is(err, constants.ErrNotFound) {
			h.bot.Send(c.Chat(), "Пользователь не найден")
			return nil
		}
		h.bot.Send(c.Chat(), "Что то пошло не так")
		logger.Error("ошибка при unbanUsername ", err)
	}
	h.bot.Send(c.Chat(), "Пользователь разбанен")
	return nil
}

func (h *Handler) btnViewBanList(c bot.Context) error {
	statistic, err := h.s.ViewBanList()
	if err != nil {
		h.bot.Send(c.Chat(), "Что то пошло не так")
		logger.Error("ошибка при получении btnViewBanList", err)
	}
	_, err = h.bot.Send(c.Chat(), statistic)
	if err != nil {
		h.bot.Send(c.Chat(), "Что то пошло не так")
		logger.Error("ошибка при отправке btnViewBanList", err)
	}
	return nil
}

func (h *Handler) newadmin(c bot.Context) error {
	username := c.Args()[0]
	url, err := h.s.NewAdmin(username)
	if err != nil {
		logger.Error("ошибка NewAdmin", err)
		h.bot.Send(c.Chat(), "Что то пошло не так")
	}
	res := fmt.Sprintf("Новый администратор должен ввести для авторизации\n /auth %s", url)
	h.bot.Send(c.Chat(), res)
	return nil
}
