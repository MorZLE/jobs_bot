package controller

import (
	"errors"
	"fmt"
	"github.com/MorZLE/jobs_bot/config"
	"github.com/MorZLE/jobs_bot/constants"
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
		4: "Прикрепите резюме в формате .pdf одной страницей",
	}

	return &Handler{
		s:                s,
		bot:              b,
		dir:              cnf.Dir,
		mQuestion:        mQestion,
		userType:         make(map[int64]string),
		userReg:          make(map[int64]model.Student),
		userQuestion:     make(map[int64]int),
		employeeCount:    make(map[int64]int),
		employeeCategory: make(map[int64]string),
	}, nil
}

type Handler struct {
	s                service.Service
	bot              *bot.Bot
	dir              string
	mQuestion        map[int]string          //Вопросы
	userReg          map[int64]model.Student //Структура студента при регистрации и id студента
	userQuestion     map[int64]int           //Номер вопроса и индекс студента
	employeeCount    map[int64]int           // Номер просматриваемого резюме и id работодателя
	employeeCategory map[int64]string        // Категория просматриваемого резюме и id работодателя
	userType         map[int64]string        // Тип пользователя работодатель или студент
}

func (h *Handler) Start() {
	h.bot.Handle("/start", h.HandlerStart)
	h.bot.Handle(&btnMainMenu, h.HandlerStart)
	h.bot.Handle(&btnEmployee, h.Employee)
	h.bot.Handle(&btnViewResumeStudents, h.btnCategorySelect)
	h.bot.Handle(&btnCreateVacance, h.btnCreateResume)
	//h.bot.Handle(&btnCategorySelect, h.btnCategorySelect)

	h.bot.Handle(&ViewResume, h.ViewRes)
	h.bot.Handle(&btnStudent, h.StudentDefault)
	h.bot.Handle(&CreateResume, h.RegStudent)

	h.bot.Handle(bot.OnText, h.Text)
	h.bot.Handle(bot.OnDocument, h.Document)
	h.bot.Handle(bot.OnPhoto, h.Document)

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

	h.bot.Start()
	log.Println("Bot started")
}

var (
	menu     = &bot.ReplyMarkup{ResizeKeyboard: true}
	selector = &bot.ReplyMarkup{}
	profile  = &bot.ReplyMarkup{}
	category = &bot.ReplyMarkup{}
	employee = &bot.ReplyMarkup{}

	// Reply buttons.

	btnEmployee = menu.Text("Я работодатель")
	btnStudent  = menu.Text("Я студент")

	btnViewResumeStudents = employee.Text("Просмотреть резюме")
	btnCreateVacance      = employee.Text("Выложить вакансию")
	btnCategorySelect     = employee.Text("Выбор категории")
	btnMainMenu           = employee.Text("Главное меню")

	CreateResume = menu.Data("Создать резюме", "createResume")

	ViewResume    = menu.Text("Профиль")
	ReviewResume  = menu.Data("Изменить резюме", "review")
	DeleteProfile = menu.Data("Удалить профиль", "deleteProfile")

	btnPrev  = selector.Data("⬅", "prev")
	btnOffer = selector.Data("Предложить работу", "Offer")
	btnNext  = selector.Data("➡", "next")
)

var (
	btnC1 = category.Data("Разработчик", "btnC1", "Разработчик")
	btnC2 = category.Data("Инфо без-ть", "btnC2", "Инфо без-ть")
	btnC3 = category.Data("Дизайнер", "btnC3", "Дизайнер")
	btnC4 = category.Data("Системный ад-р", "btnC4", "Системный ад-р")
	btnC5 = category.Data("Банковское дело", "btnC5", "Банковское дело")
	btnC6 = category.Data("Страховой агент", "btnC6", "Страховой агент")
	btnC7 = category.Data("Мечтатель", "btnC7", "Мечтатель")
)

func (h *Handler) HandlerStart(c bot.Context) error {
	menu.Reply(
		menu.Row(btnEmployee),
		menu.Row(btnStudent),
	)
	return c.Send("Привет! Я бот, который поможет тебе найти работу!", menu)
}

func (h *Handler) StudentDefault(c bot.Context) error {
	menu.Reply(
		menu.Row(btnMainMenu),
	)
	profile.Inline(
		profile.Row(CreateResume),
	)
	h.userType[c.Sender().ID] = constants.Student
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
	employee.Reply(
		employee.Row(btnViewResumeStudents),
		employee.Row(btnCreateVacance),
		employee.Row(btnMainMenu),
	)
	h.userType[c.Sender().ID] = constants.Employee
	return c.Send("Выберите действие", employee)
}

func (h *Handler) btnCreateResume(c bot.Context) error {
	h.bot.Send(c.Chat(), "Скоро будет возможность выложить вакансию :)")
	return nil
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
	urlPDF := fmt.Sprintf("src\\resume\\%d%s", user.Tgid, user.Resume)
	fmt.Println(urlPDF)
	resume := fmt.Sprintf("ФИО: %s\nГруппа: %s\nКатегория: %s\n", user.Fio, user.Group, user.Category)
	file := &bot.Photo{
		File:    bot.FromDisk(urlPDF),
		Caption: resume,
	}

	_, err = h.bot.Send(c.Chat(), file, profile)
	if err != nil {
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
	h.userQuestion[id] = 0
	delete(h.userQuestion, id)
	h.bot.Send(c.Chat(), "Профиль удален, надеюсь вы нашли работу!")
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
	h.userQuestion[id] = 1
	h.bot.Send(c.Chat(), h.mQuestion[h.userQuestion[id]])
	return nil
}

func (h *Handler) Text(c bot.Context) error {
	if h.userType[c.Sender().ID] == constants.Student {
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
		switch h.userQuestion[id] {
		case 1:
			st := model.Student{
				Tgid:     id,
				Fio:      data,
				Username: c.Sender().Username,
			}
			h.userReg[id] = st
			h.userQuestion[id]++
		case 2:
			st := h.userReg[id]
			st.Group = data
			h.userReg[id] = st
			h.userQuestion[id]++
		}
		h.RegStudent(c)
	}
	return nil
}
func (h *Handler) Document(c bot.Context) error {
	doc := c.Message().Document
	id := c.Sender().ID
	var pdfPath string
	if doc == nil {
		h.bot.Send(c.Chat(), h.mQuestion[4])
		return nil
	}
	if h.userQuestion[id] == 4 {
		if doc.FileName == "" {
			h.bot.Send(c.Chat(), h.mQuestion[4])
			return nil
		}
		if strings.HasSuffix(doc.FileName, ".pdf") {
			if strings.HasSuffix(doc.FileName, ".pdf") {
				pdfPath = fmt.Sprintf("%s\\src\\resume\\%d.pdf", h.dir, id)
				us := h.userReg[id]
				us.Resume = ".pdf"
				h.userReg[id] = us
			}
			//if strings.HasSuffix(doc.FileName, ".docx") {
			//	pdfPath = fmt.Sprintf("%s\\src\\resume\\%d.docx", h.dir, id)
			//	us := h.userReg[id]
			//	us.Resume = ".docx"
			//	h.userReg[id] = us
			//}

			err := h.s.SaveResume(h.userReg[id])
			if err != nil {
				log.Println(err)
				h.bot.Send(c.Chat(), "Что то пошло не так")
				return err
			}
			err = h.bot.Download(&doc.File, pdfPath)
			if err != nil {
				log.Println(err)
				h.bot.Send(c.Chat(), "Что то пошло не так")
				return err
			}
			h.bot.Send(c.Chat(), "Ваше резюме опубликовано")
			h.ViewRes(c)
			return nil
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
	h.userType[c.Sender().ID] = constants.Student
	if _, ok := h.userQuestion[id]; !ok {
		h.userQuestion[id] = 1
		h.bot.Send(c.Chat(), "Заполните данные о себе:")
	}
	if h.userQuestion[id] == 3 {
		category.Inline(
			category.Row(btnC1, btnC2, btnC3),
			category.Row(btnC4, btnC5, btnC6),
			category.Row(btnC7),
		)
		h.bot.Send(c.Chat(), h.mQuestion[h.userQuestion[id]], category)
		return nil
	}
	h.bot.Send(c.Chat(), h.mQuestion[h.userQuestion[id]])
	return nil
}
func (h *Handler) btnC1(c bot.Context) error {
	id := c.Sender().ID
	data := c.Data()
	switch h.userType[id] {
	case constants.Student:
		st := h.userReg[id]
		st.Category = data
		h.userReg[id] = st
		h.userQuestion[id]++
		h.RegStudent(c)
	case constants.Employee:
		h.employeeCategory[id] = data
		h.ViewResumeStudents(c, constants.Next)
	}
	return nil
}
func (h *Handler) btnCategorySelect(c bot.Context) error {
	category.Inline(
		category.Row(btnC1, btnC2, btnC3),
		category.Row(btnC4, btnC5, btnC6),
		category.Row(btnC7),
	)
	h.employeeCount[c.Sender().ID] = 0
	h.bot.Send(c.Chat(), "Выберите категорию резюме", category)
	return nil
}

func (h *Handler) ViewResumeStudents(c bot.Context, dir string) error {
	id := c.Sender().ID

	selector.Inline(
		selector.Row(btnPrev, btnOffer, btnNext),
	)
	category.Inline(
		selector.Row(btnCategorySelect),
	)

	user, count, err := h.s.GetResume(h.employeeCategory[c.Sender().ID], h.employeeCount[id], dir)
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
		log.Println(err)
	}
	h.employeeCount[id] = count
	urlPDF := fmt.Sprintf("src\\resume\\%d%s", user.Tgid, user.Resume)
	resume := fmt.Sprintf("ФИО: %s\nГруппа: %s\nКатегория: %s\n", user.Fio, user.Group, user.Category)
	file := &bot.Photo{
		File:    bot.FromDisk(urlPDF),
		Caption: resume,
	}
	h.bot.Send(c.Chat(), file, selector)
	if err != nil {
		h.bot.Send(c.Chat(), "Что то пошло не так")
	}
	return nil
}

func (h *Handler) Next(c bot.Context) error {
	id := c.Sender().ID
	h.employeeCount[id]++
	h.ViewResumeStudents(c, constants.Next)
	return nil
}
func (h *Handler) Prev(c bot.Context) error {
	id := c.Sender().ID
	if h.employeeCount[id] > 0 {
		h.employeeCount[id]--
	}
	h.ViewResumeStudents(c, constants.Prev)
	return nil
}
func (h *Handler) Offer(c bot.Context) error {
	id := c.Sender().ID
	user, _, err := h.s.GetResume(h.employeeCategory[id], h.employeeCount[id], constants.Offer)
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
