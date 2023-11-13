package controller

import (
	"github.com/MorZLE/jobs_bot/config"
	"github.com/MorZLE/jobs_bot/service"
	bot "gopkg.in/telebot.v3"
	"log"
)

func NewHandler(s service.Service, cnf *config.Config) Controller {
	b, err := bot.(cnf.BotToken)
	if err != nil {
		log.Panic(err)
	}

	return &Handler{
		s:   s,
		bot: b,
	}
}

type Handler struct {
	s   service.Service
	bot *bot.BotAPI
}

func (h *Handler) Start() {

}
