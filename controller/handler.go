package controller

import (
	"fmt"
	"github.com/MorZLE/jobs_bot/config"
	"github.com/MorZLE/jobs_bot/model"
	"github.com/MorZLE/jobs_bot/service"
	bot "gopkg.in/telebot.v3"
	"log"
	"strings"
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
		4: "Прикрепите ваше резюме в формате .pdf",
	}

	return &Handler{
		s:             s,
		bot:           b,
		dir:           cnf.Dir,
		mQuestion:     mQestion,
		userReg:       make(map[int64]model.Student),
		userQuestion:  make(map[int64]int),
		employeeCount: make(map[int64]int),
	}, nil
}

type Handler struct {
	s             service.Service
	bot           *bot.Bot
	dir           string
	mQuestion     map[int]string          //Вопросы
	userReg       map[int64]model.Student //Структура студента при регистрации и id студента
	userQuestion  map[int64]int           //Номер вопроса и индекс студента
	employeeCount map[int64]int
}

func (h *Handler) Start() {
	h.bot.Handle("/start", h.HandlerStart)
	h.bot.Handle(&btnEmployee, h.RegEmployee)
	h.bot.Handle(&ViewResume, h.ViewRes)
	h.bot.Handle(&btnStudent, h.RegStudent)
	h.bot.Handle(bot.OnText, h.Text)
	h.bot.Handle(bot.OnDocument, h.Document)
	h.bot.Handle(&btnNext, h.Next)
	h.bot.Handle(&btnPrev, h.Prev)
	h.bot.Handle(&btnOffer, h.Offer)
	h.bot.Start()
	log.Println("Bot started")
}

var (
	menu     = &bot.ReplyMarkup{ResizeKeyboard: true}
	selector = &bot.ReplyMarkup{}

	// Reply buttons.
	btnEmployee = menu.Text("Я работодатель")
	btnStudent  = menu.Text("Я студент")

	ViewResume = menu.Text("Посмотреть резюме")

	btnPrev  = selector.Data("⬅", "prev")
	btnOffer = selector.Data("Предложить работу", "prev")
	btnNext  = selector.Data("➡", "next")
)

func (h *Handler) HandlerStart(c bot.Context) error {
	menu.Reply(
		menu.Row(btnEmployee),
		menu.Row(btnStudent),
	)
	return c.Send("Привет! Я бот, который поможет тебе найти работу!", menu)
}
func (h *Handler) ViewRes(c bot.Context) error {
	user, err := h.s.Get(c.Sender().ID)
	if err != nil {
		return err
	}
	urlPDF := fmt.Sprintf("src\\pdf\\%d.pdf", user.TgId)
	resume := fmt.Sprintf("ФИО: %s\nГруппа: %s\nКатегория: %s\n", user.Fio, user.Group, user.Category)
	file := &bot.Photo{
		File:    bot.FromDisk(urlPDF),
		Caption: resume,
	}

	_, err = h.bot.Send(c.Chat(), file)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) Text(c bot.Context) error {

	data := c.Message().Text
	if data == "" {
		h.bot.Send(c.Chat(), "Введите данные")
	}
	id := c.Sender().ID
	switch h.userQuestion[id] {
	case 1:
		st := model.Student{
			TgId: id,
			Fio:  data,
		}
		h.userReg[id] = st
		h.userQuestion[id]++
	case 2:
		st := h.userReg[id]
		st.Group = data
		h.userReg[id] = st
		h.userQuestion[id]++
	case 3:
		st := h.userReg[id]
		st.Category = data
		h.userReg[id] = st
		h.userQuestion[id]++
	}
	h.RegStudent(c)
	return nil
}
func (h *Handler) Document(c bot.Context) error {
	doc := c.Message().Document
	id := c.Sender().ID
	if h.userQuestion[id] == 4 {
		if ok := strings.HasSuffix(doc.FileName, ".pdf"); ok {
			err := h.s.SaveResume(id, doc.File, h.userReg[id])
			if err != nil {
				log.Println(err)
				h.bot.Send(c.Chat(), "Что то пошло не так")
			}
			h.bot.Send(c.Chat(), "Ваше резюме опубликовано")
		} else {
			h.bot.Send(c.Chat(), "Прикрепите резюме в формате .pdf")
		}
	} else {
		h.bot.Send(c.Chat(), "Заполните данные, перед отправкой резюме")
	}

	return nil
}

func (h *Handler) RegStudent(c bot.Context) error {
	menu.Reply(
		menu.Row(ViewResume),
	)
	id := c.Sender().ID
	if _, ok := h.userQuestion[id]; !ok {
		h.userQuestion[id] = 1
		h.bot.Send(c.Chat(), "Заполните данные о себе:")
	}
	h.bot.Send(c.Chat(), h.mQuestion[h.userQuestion[id]])
	return nil
}

func (h *Handler) RegEmployee(c bot.Context) error {
	selector.Inline(
		selector.Row(btnPrev, btnOffer, btnNext),
	)
	user, err := h.s.Get(c.Sender().ID)
	if err != nil {
		return err
	}
	urlPDF := fmt.Sprintf("src\\pdf\\%d.pdf", user.TgId)
	resume := fmt.Sprintf("ФИО: %s\nГруппа: %s\nКатегория: %s\n", user.Fio, user.Group, user.Category)
	file := &bot.Photo{
		File:    bot.FromDisk(urlPDF),
		Caption: resume,
	}

	_, err = h.bot.Send(c.Chat(), file, selector)
	if err != nil {
		return err
	}

	return nil
}
func (h *Handler) Next(c bot.Context) error {
	id := c.Sender().ID
	h.employeeCount[id]++
	return nil
}
func (h *Handler) Prev(c bot.Context) error {
	id := c.Sender().ID
	if h.employeeCount[id] > 0 {
		h.employeeCount[id]--
	}
	return nil
}
func (h *Handler) Offer(c bot.Context) error {
	id := c.Sender().ID
	if h.employeeCount[id] > 0 {
		h.employeeCount[id]--
	}

	return nil
}
