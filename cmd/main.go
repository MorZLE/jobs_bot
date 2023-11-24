package main

import (
	"github.com/MorZLE/jobs_bot/config"
	"github.com/MorZLE/jobs_bot/internal/controller"
	"github.com/MorZLE/jobs_bot/internal/repository"
	"github.com/MorZLE/jobs_bot/internal/service"
	"github.com/MorZLE/jobs_bot/logger"
	"log"
)

func main() {
	cnf := config.NewConfig()
	logger.Initialize()
	st, err := repository.NewRepository(cnf)
	defer st.Close()
	if err != nil {
		log.Fatal(err)
	}

	s := service.NewService(cnf, st)
	h, err := controller.NewHandler(s, cnf)
	if err != nil {
		log.Fatal(err)
	}
	h.Start()

}
