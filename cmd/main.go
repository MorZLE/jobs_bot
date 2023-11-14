package main

import (
	"github.com/MorZLE/jobs_bot/config"
	"github.com/MorZLE/jobs_bot/controller"
	"github.com/MorZLE/jobs_bot/repository"
	"github.com/MorZLE/jobs_bot/service"
	"log"
)

func main() {
	cnf := config.NewConfig()

	st, err := repository.NewRepository(cnf)

	if err != nil {
		log.Fatal(err)
	}
	defer st.Close()

	s := service.NewService(st)
	h, err := controller.NewHandler(s, cnf)
	if err != nil {
		log.Fatal(err)
	}
	h.Start()

}
