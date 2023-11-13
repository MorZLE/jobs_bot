package controller

import (
	"github.com/MorZLE/jobs_bot/config"
	"github.com/MorZLE/jobs_bot/service"
	bot "gopkg.in/telebot.v3"
	"time"
)

func NewHandler(s service.Service, cnf *config.Config) (Controller, error) {
	pref := bot.Settings{
		Token:  cnf.BotToken,
		Poller: &bot.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := bot.NewBot(pref)
	if err != nil {
		return nil, err
	}

	return &Handler{
		s:   s,
		bot: b,
	}, nil
}

type Handler struct {
	s   service.Service
	bot *bot.Bot
}

func (h *Handler) Start() {
	h.bot.Handle("/start", func(c bot.Context) error {
		return c.Send("Привет!")
	})

	h.bot.Start()
}
