package controller

import (
	"errors"
	"fmt"
	"github.com/MorZLE/jobs_bot/config"
	"github.com/MorZLE/jobs_bot/constants"
	"github.com/MorZLE/jobs_bot/logger"
	"github.com/MorZLE/jobs_bot/model"
	"github.com/MorZLE/jobs_bot/repository"
	"github.com/MorZLE/jobs_bot/service"
	bot "gopkg.in/telebot.v3"
	"log"
	"strings"
	"sync"
	"time"
)

func NewHandler(s service.Service, cnf *config.Config) (*Handler, error) {
	pref := bot.Settings{
		Token:  cnf.BotToken,
		Poller: &bot.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := bot.NewBot(pref)
	if err != nil {
		return nil, err
	}
	mQestion := map[int]string{

		1: "Ваше ФИО",
		2: "Ваша группа",
		3: "Сфера деятельности",
		4: "Прикрепите резюме в формате .pdf одной страницей",
	}
	return &Handler{
		s:         s,
		bot:       b,
		dir:       cnf.Dir,
		mQuestion: mQestion,
		user:      make(map[int64]model.User),
	}, nil
}

type Handler struct {
	s         service.Service
	bot       *bot.Bot
	dir       string
	user      map[int64]model.User // модель пользователя
	mQuestion map[int]string
	mutex     sync.RWMutex
}

func (h *Handler) Start() {
	h.bot.Handle("/start", h.HandlerStart)
	h.bot.Handle(&btnMainMenu, h.HandlerStart)
	h.bot.Handle(&btnEmployee, h.Employee)

	h.bot.Handle(&ViewResume, h.ViewRes)
	h.bot.Handle(&btnStudent, h.StudentDefault)
	h.bot.Handle(&CreateResume, h.RegStudent)

	h.bot.Handle(bot.OnText, h.Text)
	h.bot.Handle(bot.OnDocument, h.Document)
	h.bot.Handle(bot.OnPhoto, h.Document)

	h.bot.Handle(&btnLock, h.Lock)
	h.bot.Handle(&btnNext, h.Next)
	h.bot.Handle(&btnPrev, h.Prev)
	h.bot.Handle(&btnOffer, h.Offer)
	h.bot.Handle(&ReviewResume, h.ReviewResume)
	h.bot.Handle(&DeleteProfile, h.DeleteProfile)

	h.bot.Handle(&btnC1, h.btnC1)
	h.bot.Handle(&btnC2, h.btnC1)
	h.bot.Handle(&btnC3, h.btnC1)
	h.bot.Handle(&btnC4, h.btnC1)
	h.bot.Handle(&btnC5, h.btnC1)
	h.bot.Handle(&btnC6, h.btnC1)
	h.bot.Handle(&btnC7, h.btnC1)
	h.bot.Handle(&btnC8, h.btnC1)
	h.bot.Handle(&btnC9, h.btnC1)
	h.bot.Handle(&btnC10, h.btnC1)
	h.bot.Handle(&btnC11, h.btnC1)
	h.bot.Handle(&btnC12, h.btnC1)
	h.bot.Handle(&btnC13, h.btnC1)
	h.bot.Handle(&btnC14, h.btnC1)
	h.bot.Handle(&btnC15, h.btnC1)
	h.bot.Handle(&btnC16, h.btnC1)
	h.bot.Handle(&btnC17, h.btnC1)
	h.bot.Handle(&btnC18, h.btnC1)

	h.bot.Start()
	log.Println("Bot started")
}

var (
	menu     = &bot.ReplyMarkup{ResizeKeyboard: true}
	selector = &bot.ReplyMarkup{}
	profile  = &bot.ReplyMarkup{}
	category = &bot.ReplyMarkup{}

	btnEmployee = menu.Text("Просмотреть резюме")
	btnStudent  = menu.Text("Профиль")

	btnMainMenu = menu.Text("Главное меню")

	CreateResume = menu.Data("Создать резюме", "createResume")

	ViewResume    = menu.Text("Профиль")
	ReviewResume  = menu.Data("Изменить резюме", "review")
	DeleteProfile = menu.Data("Удалить профиль", "deleteProfile")

	btnLock  = selector.Data("🔒", "lock")
	btnPrev  = selector.Data("⬅", "prev")
	btnOffer = selector.Data("Написать", "Offer")
	btnNext  = selector.Data("➡", "next")
)

var cat = repository.Category

var (
	btnC1  = category.Data(cat[0], "btnC1", cat[0])
	btnC2  = category.Data(cat[1], "btnC2", cat[1])
	btnC3  = category.Data(cat[2], "btnC3", cat[2])
	btnC4  = category.Data(cat[3], "btnC4", cat[3])
	btnC5  = category.Data(cat[4], "btnC5", cat[4])
	btnC6  = category.Data(cat[5], "btnC6", cat[5])
	btnC7  = category.Data(cat[6], "btnC7", cat[6])
	btnC8  = category.Data(cat[7], "btnC8", cat[7])
	btnC9  = category.Data(cat[8], "btnC9", cat[8])
	btnC10 = category.Data(cat[9], "btnC10", cat[9])
	btnC11 = category.Data(cat[10], "btnC11", cat[10])
	btnC12 = category.Data(cat[11], "btnC12", cat[11])
	btnC13 = category.Data(cat[12], "btnC13", cat[12])
	btnC14 = category.Data(cat[13], "btnC14", cat[13])
	btnC15 = category.Data(cat[14], "btnC15", cat[14])
	btnC16 = category.Data(cat[15], "btnC16", cat[15])
	btnC17 = category.Data(cat[16], "btnC17", cat[16])
	btnC18 = category.Data(cat[17], "btnC18", cat[17])
)

func (h *Handler) HandlerStart(c bot.Context) error {
	menu.Reply(
		menu.Row(btnEmployee),
		menu.Row(btnStudent),
	)
	m := model.User{}
	user, err := h.s.Get(c.Sender().ID)
	if err == nil {
		m = model.User{
			Student: user,
		}
	} else {
		m = model.User{
			Student: model.Student{},
		}
	}
	h.mutex.Lock()
	h.user[c.Sender().ID] = m
	h.mutex.Unlock()
	return c.Send("Привет! Я бот, который поможет тебе найти работу!", menu)
}

func (h *Handler) StudentDefault(c bot.Context) error {
	menu.Reply(
		menu.Row(btnMainMenu),
	)
	profile.Inline(
		profile.Row(CreateResume),
	)
	mUser := h.user[c.Sender().ID]
	mUser.Type = constants.Student

	h.mutex.Lock()
	h.user[c.Sender().ID] = mUser
	h.mutex.Unlock()

	user, err := h.s.Get(c.Sender().ID)
	if err != nil {
		log.Println(err)
	}
	if user != (model.Student{}) {
		return h.ViewRes(c)
	}
	return c.Send("У вас еще нет резюме \n", profile)
}

func (h *Handler) Employee(c bot.Context) error {
	h.mutex.Lock()
	mUser := h.user[c.Sender().ID]
	h.mutex.Unlock()
	mUser.Type = constants.Employee
	h.user[c.Sender().ID] = mUser

	return h.btnCategorySelect(c)
}

func (h *Handler) ViewRes(c bot.Context) error {
	profile.Inline(
		profile.Row(ReviewResume),
		profile.Row(DeleteProfile),
	)
	user, err := h.s.Get(c.Sender().ID)
	if err != nil {
		log.Println(err)
	}
	if user == (model.Student{}) {
		h.bot.Send(c.Chat(), "Профиль не найден")
		return nil
	}
	h.mutex.Lock()
	h.user[c.Sender().ID] = model.User{
		Student: user,
		Type:    constants.Student,
	}
	h.mutex.Unlock()
	urlPDF := fmt.Sprintf("src\\resume\\%s", user.Resume)
	fmt.Println(urlPDF)
	resume := fmt.Sprintf("ФИО: %s\nГруппа: %s\nКатегория: %s\n", user.Fio, user.Group, user.Category)
	file := &bot.Photo{
		File:    bot.FromDisk(urlPDF),
		Caption: resume,
	}

	_, err = h.bot.Send(c.Chat(), file, profile)
	if err != nil {
		logger.Error("ошибка ViewRes", err)
		profile.Inline(
			profile.Row(CreateResume),
		)
		return err
	}
	return nil
}

func (h *Handler) DeleteProfile(c bot.Context) error {
	id := c.Sender().ID
	user, err := h.s.Get(id)
	if err != nil {
		log.Println(err)
		h.bot.Send(c.Chat(), "Произошла ошибка")
		return err
	}
	err = h.s.Delete(id, user.Category)
	if err != nil {
		log.Println(err)
		h.bot.Send(c.Chat(), "Произошла ошибка")
		return err
	}
	h.mutex.Lock()
	mUser := h.user[c.Sender().ID]
	mUser.Nqest = 0
	mUser.Student.Status = constants.StatusDeleted
	h.user[c.Sender().ID] = mUser
	h.mutex.Unlock()
	menu.Inline(
		menu.Row(btnMainMenu),
	)
	h.bot.Send(c.Chat(), "Профиль удален, надеюсь вы нашли работу!", menu)

	return nil
}

func (h *Handler) ReviewResume(c bot.Context) error {
	id := c.Sender().ID
	user, err := h.s.Get(id)
	if err != nil {
		log.Println(err)
		h.bot.Send(c.Chat(), "Произошла ошибка")
		return err
	}
	err = h.s.Delete(id, user.Category)
	if err != nil {
		log.Println(err)
		h.bot.Send(c.Chat(), "Произошла ошибка")
		return err
	}
	h.mutex.Lock()
	mUser := h.user[c.Sender().ID]
	mUser.Student.Status = constants.StatusDeleted
	mUser.Nqest = 1
	h.user[c.Sender().ID] = mUser
	h.mutex.Unlock()
	h.bot.Send(c.Chat(), h.mQuestion[mUser.Nqest])
	return nil
}

func (h *Handler) Text(c bot.Context) error {
	id := c.Sender().ID
	mUser := h.user[id]
	if mUser.Type == constants.Student {
		data := c.Message().Text
		if data == "" {
			h.bot.Send(c.Chat(), "Введите данные")
		}
		if len(data) >= 100 {
			h.bot.Send(c.Chat(), "Слишком длинное сообщение")
			h.RegStudent(c)
			return nil
		}
		id := c.Sender().ID
		switch mUser.Nqest {
		case 1:
			st := model.Student{
				Tgid:     id,
				Fio:      data,
				Username: c.Sender().Username,
			}
			mUser.Student = st
			mUser.Nqest++
		case 2:
			mUser.Student.Group = data
			mUser.Nqest++
		}
		h.mutex.Lock()
		h.user[id] = mUser
		h.mutex.Unlock()
		h.RegStudent(c)
	}
	return nil
}
func (h *Handler) Document(c bot.Context) error {
	doc := c.Message().Document

	id := c.Sender().ID
	var pdfPath string
	mUser := h.user[id]
	if doc == nil {
		h.bot.Send(c.Chat(), h.mQuestion[4])
		return nil
	}
	if mUser.Nqest == 4 {
		if doc.FileName == "" {
			h.bot.Send(c.Chat(), h.mQuestion[4])
			return nil
		}
		if strings.HasSuffix(doc.FileName, ".pdf") {
			if strings.HasSuffix(doc.FileName, ".pdf") {

				mUser.Student.Resume = fmt.Sprintf("%d.pdf", id)
				pdfPath = fmt.Sprintf("%s\\src\\resume\\%s", h.dir, mUser.Student.Resume)
				mUser.Student.Status = constants.StatusPublished

				err := h.bot.Download(&doc.File, pdfPath)
				if err != nil {
					logger.Error("ошибка Download", err)
					h.bot.Send(c.Chat(), "Что то пошло не так")
					return err
				}

				err = h.s.SaveResume(mUser.Student)
				if err != nil {
					if errors.Is(err, constants.ErrOpenFile) {
						h.bot.Send(c.Chat(), "Не удается открыть ваше резюме, проверьте файл на целостность")
						h.bot.Send(c.Chat(), h.mQuestion[4])
						return nil
					}
					logger.Error("ошибка SaveResume", err)
					h.bot.Send(c.Chat(), "Что то пошло не так")
					return err
				}

				h.bot.Send(c.Chat(), "Ваше резюме опубликовано")
				h.ViewRes(c)
				return nil
			}
		} else {
			h.bot.Send(c.Chat(), h.mQuestion[4])
		}
	} else {
		h.bot.Send(c.Chat(), "Заполните данные, перед отправкой резюме")
	}
	return nil
}

func (h *Handler) RegStudent(c bot.Context) error {
	id := c.Sender().ID
	mUser := h.user[id]
	if mUser.Student.Status == constants.StatusPublished {
		return nil
	}
	mUser.Type = constants.Student
	if mUser.Nqest == 0 {
		mUser.Nqest = 1
		h.mutex.Lock()
		h.user[id] = mUser
		h.mutex.Unlock()
		h.bot.Send(c.Chat(), "Заполните данные о себе:")
	}
	if mUser.Nqest == 3 {
		category := GetCategory()
		h.bot.Send(c.Chat(), h.mQuestion[mUser.Nqest], category)
		return nil
	}
	h.bot.Send(c.Chat(), h.mQuestion[mUser.Nqest])
	return nil
}
func (h *Handler) btnC1(c bot.Context) error {
	id := c.Sender().ID
	data := c.Data()
	mUser := h.user[id]
	switch mUser.Type {
	case constants.Student:
		if mUser.Student.Status == constants.StatusPublished {
			return nil
		}
		if mUser.Nqest != 3 {
			return nil
		}
		mUser.Student.Category = data
		mUser.Nqest++
		h.mutex.Lock()
		h.user[id] = mUser
		h.mutex.Unlock()
		h.RegStudent(c)
		h.bot.Delete(c.Message())
	case constants.Employee:
		mUser.EmployeeCategory = data
		mUser.EmployeeCount = 0
		mUser.EmployeeSetCategory = true
		h.mutex.Lock()
		h.user[id] = mUser
		h.mutex.Unlock()
		h.ViewResumeStudents(c, constants.Next)
	}
	return nil
}
func GetCategory() *bot.ReplyMarkup {
	category.Inline(
		category.Row(btnC1),
		category.Row(btnC2),
		category.Row(btnC3),
		category.Row(btnC4),
		category.Row(btnC5),
		category.Row(btnC6),
		category.Row(btnC7),
		category.Row(btnC8),
		category.Row(btnC9),
		category.Row(btnC10),
		category.Row(btnC11),
		category.Row(btnC12),
		category.Row(btnC13),
		category.Row(btnC14),
		category.Row(btnC15),
		category.Row(btnC16),
		category.Row(btnC17),
		category.Row(btnC18),
	)
	return category
}
func (h *Handler) btnCategorySelect(c bot.Context) error {
	category := GetCategory()
	id := c.Sender().ID
	mUser := h.user[id]
	mUser.EmployeeCount = 0
	mUser.EmployeeSetCategory = true
	h.mutex.Lock()
	h.user[id] = mUser
	h.mutex.Unlock()
	_, err := h.bot.Send(c.Chat(), "Выберите категорию резюме", category)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func (h *Handler) ViewResumeStudents(c bot.Context, dir string) error {
	id := c.Sender().ID

	mUser := h.user[id]

	user, count, err := h.s.GetResume(mUser.EmployeeCategory, mUser.EmployeeCount, dir)
	if count > 0 {
		selector.Inline(
			selector.Row(btnPrev, btnOffer, btnNext),
		)
	} else if count == 0 {
		selector.Inline(
			selector.Row(btnLock, btnOffer, btnNext),
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
				selector.Row(btnPrev, btnOffer, btnLock),
			)
			if count == 0 {
				selector.Inline(
					selector.Row(btnLock, btnOffer, btnLock),
				)
			}
		} else {
			h.bot.Send(c.Chat(), "Что то пошло не так")
			log.Println(err)
		}
	}

	mUser.EmployeeCount = count
	urlPDF := fmt.Sprintf("src\\resume\\%s", user.Resume)
	resume := fmt.Sprintf("ФИО: %s\nГруппа: %s\nКатегория: %s\n", user.Fio, user.Group, user.Category)
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

func (h *Handler) Next(c bot.Context) error {
	id := c.Sender().ID

	mUser := h.user[id]
	mUser.EmployeeCount++
	h.mutex.Lock()
	h.user[id] = mUser
	h.mutex.Unlock()
	h.ViewResumeStudents(c, constants.Next)
	return nil
}
func (h *Handler) Prev(c bot.Context) error {
	id := c.Sender().ID
	mUser := h.user[id]
	if mUser.EmployeeCount > 0 {
		mUser.EmployeeCount--
	}
	h.mutex.Lock()
	h.user[id] = mUser
	h.mutex.Unlock()
	h.ViewResumeStudents(c, constants.Prev)
	return nil
}
func (h *Handler) Offer(c bot.Context) error {
	id := c.Sender().ID
	mUser := h.user[id]
	user, _, err := h.s.GetResume(mUser.EmployeeCategory, mUser.EmployeeCount, constants.Offer)
	if err != nil {
		if errors.Is(err, constants.ErrNotCategory) {
			h.bot.Send(c.Chat(), "Категория не найдена")
			return nil
		}
		if errors.Is(err, constants.ErrNotFound) {
			h.bot.Send(c.Chat(), "Профиль не найден")
			return nil
		}
		log.Println(err)
	}
	h.bot.Send(c.Chat(), fmt.Sprintf("Надеюсь вам понравится этот кандидат, его профиль: @%s", user.Username))
	return nil
}

func (h *Handler) Lock(c bot.Context) error {
	return nil
}
